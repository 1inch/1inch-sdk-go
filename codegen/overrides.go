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
}

// ref builds a $ref schema node.
func ref(name string) spec {
	return spec{"$ref": "#/components/schemas/" + name}
}

// specOverrides lists the known type bugs in each upstream spec, keyed by spec
// file name. See schemaPatch for semantics.
var specOverrides = map[string][]schemaPatch{
	"tokens-openapi.json": {
		// The tokens API returns tag objects {provider, value}, not strings.
		{path: "components.schemas.ProviderTokenDto.properties.tags.items", value: ref("TagDto")},
		{path: "components.schemas.TokenInfoDto.properties.tags.items", value: ref("TagDto")},
	},
}

// applyOverrides applies the spec's patch list, if any.
func applyOverrides(specName string, doc spec) error {
	for _, patch := range specOverrides[specName] {
		node, err := resolvePath(doc, patch.path)
		if err != nil {
			return fmt.Errorf("override %q: %w", patch.path, err)
		}
		for k := range node {
			delete(node, k)
		}
		for k, v := range patch.value {
			node[k] = v
		}
	}
	return nil
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
