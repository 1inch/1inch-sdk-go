// Package codegen generates the SDK's *_types.gen.go files from the OpenAPI
// specs in codegen/openapi.
//
// The upstream specs contain known bugs (wrong types, malformed schemas,
// machine-generated operation ids). Rather than hand-maintaining corrected
// copies of the generated types, this package applies a fixed set of
// transforms to each spec in memory — the committed spec files are never
// modified — and feeds the corrected spec to oapi-codegen.
//
// The behavior of every transform is pinned by the characterization tests in
// codegen_test.go; the committed *_types.gen.go files must be reproducible
// byte-for-byte from the committed specs.
package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

// spec is a decoded OpenAPI document. Numbers are decoded as json.Number so
// values above 2^53 survive a decode/encode round trip.
type spec = map[string]any

// chainIDKey matches parameter and property names that hold chain ids, e.g.
// chainId, srcChain, dstChainId, chain_id, chainIds.
var chainIDKey = regexp.MustCompile(`(?i)^(src|dst|from|to)?_?chain_?(id)?s?$`)

// transform is one in-memory spec correction. Transforms run in order; each
// receives the output of the previous one.
type transform func(doc spec) error

// specTransforms is the transform chain applied to every spec before type
// generation, in order.
var specTransforms = []transform{
	fixIncorrectNumberArrays,
	removeNullAnyOfVariants,
	collapseAnyOfParameters,
	simplifyAllOfRefs,
	fixChainIDNumberTypes,
	addSkipOptionalPointer,
}

// walkAny visits every map and slice in the document tree, calling visit on
// each map after its children have been visited.
func walkAny(node any, visit func(m spec)) {
	switch v := node.(type) {
	case spec:
		for _, child := range v {
			walkAny(child, visit)
		}
		visit(v)
	case []any:
		for _, child := range v {
			walkAny(child, visit)
		}
	}
}

// forEachOperation calls fn for every operation object under .paths.
func forEachOperation(doc spec, fn func(op spec)) {
	paths, _ := doc["paths"].(spec)
	for _, pathItem := range paths {
		item, ok := pathItem.(spec)
		if !ok {
			continue
		}
		for _, opVal := range item {
			if op, ok := opVal.(spec); ok {
				fn(op)
			}
		}
	}
}

// forEachParameter calls fn for every parameter object of every operation.
func forEachParameter(doc spec, fn func(param spec)) {
	forEachOperation(doc, func(op spec) {
		params, _ := op["parameters"].([]any)
		for _, p := range params {
			if param, ok := p.(spec); ok {
				fn(param)
			}
		}
	})
}

// fixIncorrectNumberArrays rewrites the invalid schema type "number[]" (used by
// some upstream specs) into a proper array-of-number schema.
func fixIncorrectNumberArrays(doc spec) error {
	walkAny(doc, func(m spec) {
		if m["type"] == "number[]" {
			m["type"] = "array"
			m["items"] = spec{"type": "number"}
		}
	})
	return nil
}

// renameOperationIDs replaces machine-generated operation ids with the
// friendly names from the mapping, which oapi-codegen uses to name the
// generated params types.
func renameOperationIDs(doc spec, mapping map[string]string) error {
	forEachOperation(doc, func(op spec) {
		if id, ok := op["operationId"].(string); ok {
			if mapped, ok := mapping[id]; ok {
				op["operationId"] = mapped
			}
		}
	})
	return nil
}

// removeNullAnyOfVariants drops `{"type": "null"}` variants from every anyOf
// in the document, so a nullable field generates as its non-null type.
func removeNullAnyOfVariants(doc spec) error {
	walkAny(doc, func(m spec) {
		variants, ok := m["anyOf"].([]any)
		if !ok {
			return
		}
		kept := make([]any, 0, len(variants))
		for _, v := range variants {
			if vm, ok := v.(spec); ok && vm["type"] == "null" {
				continue
			}
			kept = append(kept, v)
		}
		m["anyOf"] = kept
	})
	return nil
}

// collapseAnyOfParameters replaces an anyOf parameter schema with its first
// variant (plus the schema's description and title), because oapi-codegen
// would otherwise generate an unusable union type for a simple query
// parameter.
func collapseAnyOfParameters(doc spec) error {
	forEachParameter(doc, func(param spec) {
		schema, ok := param["schema"].(spec)
		if !ok {
			return
		}
		variants, ok := schema["anyOf"].([]any)
		if !ok || len(variants) == 0 {
			return
		}
		first, ok := variants[0].(spec)
		if !ok {
			return
		}
		collapsed := spec{}
		for k, v := range first {
			collapsed[k] = v
		}
		if d, ok := schema["description"]; ok {
			collapsed["description"] = d
		} else {
			collapsed["description"] = nil
		}
		if t, ok := schema["title"]; ok {
			collapsed["title"] = t
		} else {
			collapsed["title"] = nil
		}
		param["schema"] = collapsed
	})
	return nil
}

// simplifyAllOfRef collapses an allOf that wraps a single $ref into the bare
// ref. Non-ref allOf members are dropped. Applied to parameter schemas,
// request bodies, component schemas, and component schema properties —
// matching where the specs actually use the pattern.
func simplifyAllOfRef(schema spec) {
	members, ok := schema["allOf"].([]any)
	if !ok {
		return
	}
	refs := make([]any, 0, len(members))
	for _, m := range members {
		if mm, ok := m.(spec); ok {
			if _, isRef := mm["$ref"]; isRef {
				refs = append(refs, mm)
			}
		}
	}
	if len(refs) == 1 {
		ref := refs[0].(spec)
		for k := range schema {
			delete(schema, k)
		}
		for k, v := range ref {
			schema[k] = v
		}
		return
	}
	// Faithful to the original jq transform: an allOf that does not reduce to
	// exactly one ref is dropped entirely.
	delete(schema, "allOf")
}

func simplifyAllOfRefs(doc spec) error {
	forEachParameter(doc, func(param spec) {
		if schema, ok := param["schema"].(spec); ok {
			simplifyAllOfRef(schema)
		}
	})
	forEachOperation(doc, func(op spec) {
		body, ok := op["requestBody"].(spec)
		if !ok {
			return
		}
		content, _ := body["content"].(spec)
		jsonContent, _ := content["application/json"].(spec)
		if schema, ok := jsonContent["schema"].(spec); ok {
			simplifyAllOfRef(schema)
		}
	})
	components, _ := doc["components"].(spec)
	schemas, _ := components["schemas"].(spec)
	for _, s := range schemas {
		schema, ok := s.(spec)
		if !ok {
			continue
		}
		simplifyAllOfRef(schema)
		props, _ := schema["properties"].(spec)
		for _, p := range props {
			if prop, ok := p.(spec); ok {
				simplifyAllOfRef(prop)
			}
		}
	}
	return nil
}

// fixChainIDNumberTypes converts chain-id parameters and properties declared
// as "number" to "integer". Chain ids are integers; a float32 cannot represent
// ids above 2^24 (Aurora, 1313161554, silently rounds to 1313161600).
func fixChainIDNumberTypes(doc spec) error {
	walkAny(doc, func(m spec) {
		if props, ok := m["properties"].(spec); ok {
			for key, p := range props {
				prop, ok := p.(spec)
				if !ok || !chainIDKey.MatchString(key) {
					continue
				}
				if prop["type"] == "number" {
					prop["type"] = "integer"
				}
			}
		}
		name, _ := m["name"].(string)
		if name == "" || !chainIDKey.MatchString(name) {
			return
		}
		if schema, ok := m["schema"].(spec); ok && schema["type"] == "number" {
			schema["type"] = "integer"
		}
	})
	return nil
}

// addSkipOptionalPointer annotates schemas with
// x-go-type-skip-optional-pointer so oapi-codegen generates optional fields as
// values instead of pointers. Applied to optional operation parameters,
// request-body schemas and their properties, and component schema properties;
// only schemas with a type or oneOf are annotated (bare $refs keep their
// pointer).
func addSkipOptionalPointer(doc spec) error {
	annotate := func(schema spec) {
		_, hasType := schema["type"]
		_, hasOneOf := schema["oneOf"]
		if !hasType && !hasOneOf {
			return
		}
		if schema["x-go-type-skip-optional-pointer"] == true {
			return
		}
		schema["x-go-type-skip-optional-pointer"] = true
	}

	forEachParameter(doc, func(param spec) {
		if param["required"] != false {
			return
		}
		if schema, ok := param["schema"].(spec); ok {
			annotate(schema)
		}
	})
	forEachOperation(doc, func(op spec) {
		body, ok := op["requestBody"].(spec)
		if !ok {
			return
		}
		content, _ := body["content"].(spec)
		jsonContent, _ := content["application/json"].(spec)
		schema, ok := jsonContent["schema"].(spec)
		if !ok {
			return
		}
		annotate(schema)
		props, _ := schema["properties"].(spec)
		for _, p := range props {
			if prop, ok := p.(spec); ok {
				annotate(prop)
			}
		}
	})
	components, _ := doc["components"].(spec)
	schemas, _ := components["schemas"].(spec)
	for _, s := range schemas {
		schema, ok := s.(spec)
		if !ok {
			continue
		}
		props, _ := schema["properties"].(spec)
		for _, p := range props {
			if prop, ok := p.(spec); ok {
				annotate(prop)
			}
		}
	}
	return nil
}

// applyTransforms decodes a spec, runs the operation-id mapping and the full
// transform chain, and returns the corrected document re-encoded as JSON.
func applyTransforms(raw []byte, operationIDs map[string]string) ([]byte, error) {
	doc, err := decodeSpec(raw)
	if err != nil {
		return nil, err
	}
	if err := renameOperationIDs(doc, operationIDs); err != nil {
		return nil, err
	}
	for _, t := range specTransforms {
		if err := t(doc); err != nil {
			return nil, err
		}
	}
	return json.Marshal(doc)
}

func decodeSpec(raw []byte) (spec, error) {
	var doc spec
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("invalid spec JSON: %w", err)
	}
	return doc, nil
}
