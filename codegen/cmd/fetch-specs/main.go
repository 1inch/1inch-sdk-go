// Command fetch-specs refreshes the local OpenAPI spec copies from the 1inch
// Dev Portal. Run from the repository root:
//
//	DEV_PORTAL_TOKEN=... go run ./codegen/cmd/fetch-specs
//
// or via `make codegen-fetch-specs`. Review the spec diff, then run
// `make codegen-types` to regenerate.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/1inch/1inch-sdk-go/v4/codegen"
)

func main() {
	specsDir := flag.String("specs", "codegen/openapi", "directory to write *-openapi.json spec files to")
	mappingFile := flag.String("mapping", "codegen/mapping.json", "operation-id mapping file")
	token := flag.String("token", os.Getenv("DEV_PORTAL_TOKEN"), "Dev Portal API token (defaults to DEV_PORTAL_TOKEN)")
	flag.Parse()

	err := codegen.FetchSpecs(codegen.FetchOptions{
		SpecsDir:    *specsDir,
		MappingFile: *mappingFile,
		Token:       *token,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch-specs:", err)
		os.Exit(1)
	}
}
