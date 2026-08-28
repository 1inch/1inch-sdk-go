package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// generatorModule pins the exact oapi-codegen build used for generation. It is
// invoked hermetically via `go run module@version`, so no pre-installed binary
// is required and the version cannot drift between machines.
const generatorModule = "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0"

const specSuffix = "-openapi.json"

// formTagPattern rewrites the `form:"..."` struct tags oapi-codegen emits on
// params types into the `url:"..."` tags used by go-querystring.
var formTagPattern = regexp.MustCompile("`form:\"([^\"]+)\"([^`]*)`")

// Options configures a generation run.
type Options struct {
	// SpecsDir is the directory containing the *-openapi.json spec files.
	SpecsDir string
	// MappingFile is the JSON file mapping upstream operation ids to friendly
	// names.
	MappingFile string
	// OutputDir is the sdk-clients directory that receives
	// <package>/<spec>_types.gen.go files.
	OutputDir string
}

// Generate runs the full type-generation pipeline: it validates spec file
// names, applies the transform chain to each spec in memory (the spec files on
// disk are never modified), invokes the pinned oapi-codegen on the corrected
// spec, rewrites form tags to url tags, and writes the result into the output
// directory.
func Generate(opts Options) error {
	mapping, err := loadOperationIDMapping(opts.MappingFile)
	if err != nil {
		return err
	}

	specs, err := listSpecFiles(opts.SpecsDir)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("no %s files found in %s", specSuffix, opts.SpecsDir)
	}

	for _, name := range specs {
		if err := generateOne(opts, name, mapping); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

// listSpecFiles returns the spec file names in the directory, sorted, and
// rejects any file that does not follow the <name>-openapi.json convention.
func listSpecFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var specs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), specSuffix) {
			return nil, fmt.Errorf("file %q does not match the expected naming schema of *%s", e.Name(), specSuffix)
		}
		specs = append(specs, e.Name())
	}
	sort.Strings(specs)
	return specs, nil
}

// packageName derives the Go package for a spec file name: everything before
// the first underscore. fusionplus_orders-openapi.json → package fusionplus,
// file fusionplus_orders_types.gen.go.
func packageName(specName string) string {
	base := strings.TrimSuffix(specName, specSuffix)
	if i := strings.Index(base, "_"); i >= 0 {
		return base[:i]
	}
	return base
}

func outputFileName(specName string) string {
	return strings.TrimSuffix(specName, specSuffix) + "_types.gen.go"
}

func generateOne(opts Options, specName string, mapping map[string]string) error {
	raw, err := os.ReadFile(filepath.Join(opts.SpecsDir, specName))
	if err != nil {
		return err
	}

	corrected, err := applyTransforms(raw, mapping)
	if err != nil {
		return err
	}

	generated, err := runGenerator(packageName(specName), corrected)
	if err != nil {
		return err
	}

	generated = formTagPattern.ReplaceAll(generated, []byte("`url:\"$1\"$2`"))

	outDir := filepath.Join(opts.OutputDir, packageName(specName))
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, outputFileName(specName)), generated, 0o644)
}

// runGenerator writes the corrected spec to a temp file and runs the pinned
// oapi-codegen over it, returning the generated Go source.
func runGenerator(pkg string, correctedSpec []byte) ([]byte, error) {
	tmp, err := os.CreateTemp("", "spec-*.json")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(correctedSpec); err != nil {
		tmp.Close()
		return nil, err
	}
	if err := tmp.Close(); err != nil {
		return nil, err
	}

	cmd := exec.Command("go", "run", generatorModule, "-generate", "types", "-package", pkg, tmp.Name())
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("oapi-codegen failed: %w\n%s", err, stderr.String())
	}
	return out, nil
}

func loadOperationIDMapping(path string) (map[string]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var mapping map[string]string
	if err := json.Unmarshal(raw, &mapping); err != nil {
		return nil, fmt.Errorf("invalid operation-id mapping file %s: %w", path, err)
	}
	return mapping, nil
}
