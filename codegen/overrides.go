package codegen

import (
	"fmt"
	"strings"
)

// schemaPatch replaces the entire content of one node in a spec document.
// Patches are the declarative record of every place the upstream specs are
// wrong: instead of hand-maintaining corrected copies of generated types
// ("*Fixed" structs), the spec is corrected here and the generated types come
// out right.
//
// The path is dot-separated map navigation from the document root. Applying a
// patch fails loudly if the path does not resolve, so a refreshed upstream
// spec that moves or fixes a field forces this table to be revisited.
type schemaPatch struct {
	path  string
	value spec
	// merge adds/overwrites the value's keys on the node instead of replacing
	// the node's entire content. Used to add properties a spec omits.
	merge bool
}

// ref builds a $ref schema node.
func ref(name string) spec {
	return spec{"$ref": "#/components/schemas/" + name}
}

// bigIntType is an x-go-type schema generating *big.Int, for amounts that
// overflow float64.
func bigIntType() spec {
	return spec{
		"type":             "number",
		"x-go-type":        "*big.Int",
		"x-go-type-import": spec{"path": "math/big"},
	}
}

// specOverrides lists the known type bugs in each upstream spec, keyed by spec
// file name. See schemaPatch for semantics.
var specOverrides = map[string][]schemaPatch{
	"tokens-openapi.json": {
		// The tokens API returns tag objects {provider, value}, not strings.
		{path: "components.schemas.ProviderTokenDto.properties.tags.items", value: ref("TagDto")},
		{path: "components.schemas.TokenInfoDto.properties.tags.items", value: ref("TagDto")},
	},
	"fusionplus_quoter-openapi.json": {
		// quoteId is declared as object but the API returns a string.
		{path: "components.schemas.GetQuoteOutput.properties.quoteId", value: spec{
			"type":        "string",
			"description": "Current generated quote id, should be passed with order",
		}},
	},
	"fusion_quoter-openapi.json": {
		// quoteId is declared as object but the API returns a string.
		{path: "components.schemas.GetQuoteOutput.properties.quoteId", value: spec{
			"type":        "string",
			"description": "Current generated quote id, should be passed with order",
		}},
		// exclusiveResolver is declared as object but the API returns the
		// resolver address as a string.
		{path: "components.schemas.PresetClass.properties.exclusiveResolver", value: spec{"type": "string"}},
		// The quoter also returns surplusFee and marketAmount, which the spec
		// omits.
		{path: "components.schemas.GetQuoteOutput.properties", merge: true, value: spec{
			"surplusFee":   spec{"type": "number"},
			"marketAmount": spec{"type": "string"},
		}},
	},
	"fusionplus_orders-openapi.json": {
		// Token USD prices are declared as objects but the API returns
		// decimal strings.
		{path: "components.schemas.GetOrderFillsByHashOutput.properties.dstTokenPriceUsd", value: spec{"type": "string"}},
		{path: "components.schemas.GetOrderFillsByHashOutput.properties.srcTokenPriceUsd", value: spec{"type": "string"}},
		// points is declared as a single object but the API returns an array.
		{path: "components.schemas.GetOrderFillsByHashOutput.properties.points", value: spec{
			"type":  "array",
			"items": ref("AuctionPointOutput"),
		}},
	},
}

// paramOverrides corrects the schema of every operation parameter with the
// given name, keyed by spec file name. Used when an upstream spec declares the
// same wrongly-typed parameter across several endpoints.
var paramOverrides = map[string]map[string]spec{
	"fusionplus_quoter-openapi.json": {
		// Amounts are big-integer wei strings; float32 would corrupt them.
		"amount": {"type": "string"},
		// Fee is in bps; *big.Int matches the SDK-wide amount handling.
		"fee": bigIntType(),
		// isPermit2 is a flag, not a string.
		"isPermit2": {"type": "boolean"},
	},
	"fusion_quoter-openapi.json": {
		// Amounts are big-integer wei strings; float32 would corrupt them.
		"amount": {"type": "string"},
	},
}

// paramAdditions appends parameters the upstream spec omits, keyed by spec
// file name and raw (pre-mapping) operation id. Appending keeps the new fields
// at the end of the generated params struct.
var paramAdditions = map[string]map[string][]spec{
	"fusion_quoter-openapi.json": {
		"QuoterController_getQuote": {
			surplusParam(),
		},
		"QuoterController_getQuoteWithCustomPresets": {
			spec{
				"name":     "isLedgerLive",
				"in":       "query",
				"required": true,
				"schema":   spec{"type": "boolean"},
			},
			surplusParam(),
		},
	},
}

// surplusParam is the surplus flag the quoter accepts but the spec omits.
func surplusParam() spec {
	return spec{
		"name":     "surplus",
		"in":       "query",
		"required": false,
		"schema":   spec{"type": "boolean"},
	}
}

// applyOverrides applies the spec's schema patches and parameter overrides,
// if any.
func applyOverrides(specName string, doc spec) error {
	for _, patch := range specOverrides[specName] {
		node, err := resolvePath(doc, patch.path)
		if err != nil {
			return fmt.Errorf("override %q: %w", patch.path, err)
		}
		if patch.merge {
			for k, v := range patch.value {
				node[k] = v
			}
		} else {
			replaceContent(node, patch.value)
		}
	}

	if err := addParameters(specName, doc); err != nil {
		return err
	}

	params := paramOverrides[specName]
	if len(params) == 0 {
		return nil
	}
	applied := make(map[string]bool, len(params))
	forEachParameter(doc, func(param spec) {
		name, _ := param["name"].(string)
		value, ok := params[name]
		if !ok {
			return
		}
		schema, ok := param["schema"].(spec)
		if !ok {
			return
		}
		replaceContent(schema, value)
		applied[name] = true
	})
	for name := range params {
		if !applied[name] {
			return fmt.Errorf("parameter override %q matched no parameter", name)
		}
	}
	return nil
}

// addParameters appends the spec's parameter additions, if any. A parameter is
// only added if the operation does not already declare it, so a refreshed
// upstream spec that gains the parameter wins over the addition.
func addParameters(specName string, doc spec) error {
	additions := paramAdditions[specName]
	if len(additions) == 0 {
		return nil
	}
	applied := make(map[string]bool, len(additions))
	forEachOperation(doc, func(op spec) {
		id, _ := op["operationId"].(string)
		add, ok := additions[id]
		if !ok {
			return
		}
		params, _ := op["parameters"].([]any)
		declared := make(map[string]bool, len(params))
		for _, p := range params {
			if pm, ok := p.(spec); ok {
				if name, ok := pm["name"].(string); ok {
					declared[name] = true
				}
			}
		}
		for _, param := range add {
			if name, _ := param["name"].(string); !declared[name] {
				params = append(params, param)
			}
		}
		op["parameters"] = params
		applied[id] = true
	})
	for id := range additions {
		if !applied[id] {
			return fmt.Errorf("parameter addition for operation %q matched no operation", id)
		}
	}
	return nil
}

// replaceContent replaces a node's entire content with the given value.
func replaceContent(node, value spec) {
	for k := range node {
		delete(node, k)
	}
	for k, v := range value {
		node[k] = v
	}
}

func resolvePath(doc spec, path string) (spec, error) {
	node := doc
	for _, key := range strings.Split(path, ".") {
		child, ok := node[key].(spec)
		if !ok {
			return nil, fmt.Errorf("path element %q not found or not an object", key)
		}
		node = child
	}
	return node, nil
}
