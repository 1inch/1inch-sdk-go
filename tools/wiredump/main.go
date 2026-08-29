// Command wiredump prints the SDK's wire surface: every exported struct field
// of every package under the given root, with its Go type and struct tag, as
// deterministic JSON.
//
// Struct tags are invisible to compile-time API checks (gorelease/apidiff): a
// changed json or url tag compiles fine but silently changes what goes over
// the wire. Diffing two wiredump outputs (e.g. a PR against its base) surfaces
// every field addition, removal, retyping, and retagging in one place:
//
//	go run ./tools/wiredump . > new.json
//	go run ./tools/wiredump /path/to/base-checkout > old.json
//	diff old.json new.json
//
// The dump is intentionally syntactic (types as source text, from the AST):
// the goal is "did anything change", with classification left to the reader.
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// field is one exported struct field's wire-relevant identity.
type field struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  string `json:"tag,omitempty"`
}

// dumpDirs are the roots (relative to the repo root) whose exported structs
// form the SDK's public wire surface.
var dumpDirs = []string{"sdk-clients", "common", "constants"}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	// package path -> type name -> fields
	surface := map[string]map[string][]field{}

	for _, dir := range dumpDirs {
		base := filepath.Join(root, dir)
		err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return err
			}
			if d.Name() == "examples" || d.Name() == "testdata" {
				return filepath.SkipDir
			}
			pkgTypes, err := dumpPackage(path)
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			if len(pkgTypes) > 0 {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				surface[filepath.ToSlash(rel)] = pkgTypes
			}
			return nil
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "wiredump:", err)
			os.Exit(1)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(surface); err != nil {
		fmt.Fprintln(os.Stderr, "wiredump:", err)
		os.Exit(1)
	}
}

// dumpPackage parses the non-test Go files of one directory and returns its
// exported struct types.
func dumpPackage(dir string) (map[string][]field, error) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	types := map[string][]field{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		{
			ast.Inspect(f, func(n ast.Node) bool {
				ts, ok := n.(*ast.TypeSpec)
				if !ok || !ts.Name.IsExported() {
					return true
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return true
				}
				fields := []field{}
				for _, fld := range st.Fields.List {
					typeText := nodeText(fset, dir, fld.Type)
					tag := ""
					if fld.Tag != nil {
						tag = strings.Trim(fld.Tag.Value, "`")
					}
					if len(fld.Names) == 0 {
						// embedded field
						fields = append(fields, field{Name: typeText, Type: "(embedded)", Tag: tag})
						continue
					}
					for _, name := range fld.Names {
						if !name.IsExported() {
							continue
						}
						fields = append(fields, field{Name: name.Name, Type: typeText, Tag: tag})
					}
				}
				sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
				types[ts.Name.Name] = fields
				return true
			})
		}
	}
	return types, nil
}

// nodeText renders an AST expression back to its source text.
func nodeText(fset *token.FileSet, dir string, n ast.Node) string {
	start := fset.Position(n.Pos())
	end := fset.Position(n.End())
	src, err := os.ReadFile(start.Filename)
	if err != nil {
		return "?"
	}
	return string(src[start.Offset:end.Offset])
}
