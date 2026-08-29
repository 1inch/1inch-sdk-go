// Command generate-types regenerates all *_types.gen.go files from the
// OpenAPI specs in codegen/openapi. Run from the repository root:
//
//	go run ./codegen/cmd/generate-types
//
// or via `make codegen-types`.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/1inch/1inch-sdk-go/v5/codegen"
)

func main() {
	specsDir := flag.String("specs", "codegen/openapi", "directory containing *-openapi.json spec files")
	mappingFile := flag.String("mapping", "codegen/mapping.json", "operation-id mapping file")
	outputDir := flag.String("out", "sdk-clients", "output directory for generated types")
	flag.Parse()

	err := codegen.Generate(codegen.Options{
		SpecsDir:    *specsDir,
		MappingFile: *mappingFile,
		OutputDir:   *outputDir,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate-types:", err)
		os.Exit(1)
	}
}
