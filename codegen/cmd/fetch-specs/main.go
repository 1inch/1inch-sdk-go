// Command fetch-specs refreshes the local OpenAPI spec copies from the 1inch
// Dev Portal and records their provenance in codegen/specs.lock.json. Run from
// the repository root:
//
//	DEV_PORTAL_TOKEN=... go run ./codegen/cmd/fetch-specs
//
// or via `make codegen-fetch-specs`. Review the spec diff, then run
// `make codegen-types` to regenerate.
//
// The -seed flag rebuilds the lock file from the spec copies already on disk
// without any network access (bootstrap/repair only).
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
	lockFile := flag.String("lock", "codegen/specs.lock.json", "spec provenance lock file")
	token := flag.String("token", os.Getenv("DEV_PORTAL_TOKEN"), "Dev Portal API token (defaults to DEV_PORTAL_TOKEN)")
	seed := flag.Bool("seed", false, "rebuild the lock file from the local spec copies instead of fetching")
	flag.Parse()

	opts := codegen.FetchOptions{
		SpecsDir:    *specsDir,
		MappingFile: *mappingFile,
		LockFile:    *lockFile,
		Token:       *token,
	}

	var err error
	if *seed {
		err = codegen.SeedLock(opts)
	} else {
		err = codegen.FetchSpecs(opts)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch-specs:", err)
		os.Exit(1)
	}
}
