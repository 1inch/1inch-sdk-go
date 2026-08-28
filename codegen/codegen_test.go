// Package codegen contains characterization tests for the generate_types.sh
// pipeline. They pin down the pipeline's current observable behavior so it can
// be safely restructured: any reimplementation must keep these tests green (or
// change them in a reviewed commit alongside the regenerated files).
//
// The tests run the real script in a temporary copy of the repo layout, so they
// require bash, jq, and sed on PATH (all present on CI). The script installs
// the pinned oapi-codegen via `go install`, which needs the module cache or
// network on first run.
package codegen

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// specFiles is the complete set of OpenAPI specs the pipeline consumes. A new
// spec must be added here (and to generatedFiles) so the reproducibility test
// covers it.
var specFiles = []string{
	"aggregation-openapi.json",
	"balances-openapi.json",
	"fusion_orders-openapi.json",
	"fusion_quoter-openapi.json",
	"fusion_relayer-openapi.json",
	"fusionplus_orders-openapi.json",
	"fusionplus_quoter-openapi.json",
	"fusionplus_relayer-openapi.json",
	"gasprices-openapi.json",
	"history-openapi.json",
	"nft-openapi.json",
	"orderbook-openapi.json",
	"portfolio-openapi.json",
	"spotprices-openapi.json",
	"tokens-openapi.json",
	"traces-openapi.json",
	"txbroadcast-openapi.json",
	"web3-openapi.json",
}

// generatedFiles is the complete set of files the pipeline produces, relative
// to the sdk-clients directory. It also pins the package-placement rule: a spec
// named <pkg>_<suffix>-openapi.json generates <pkg>/<pkg>_<suffix>_types.gen.go.
var generatedFiles = []string{
	"aggregation/aggregation_types.gen.go",
	"balances/balances_types.gen.go",
	"fusion/fusion_orders_types.gen.go",
	"fusion/fusion_quoter_types.gen.go",
	"fusion/fusion_relayer_types.gen.go",
	"fusionplus/fusionplus_orders_types.gen.go",
	"fusionplus/fusionplus_quoter_types.gen.go",
	"fusionplus/fusionplus_relayer_types.gen.go",
	"gasprices/gasprices_types.gen.go",
	"history/history_types.gen.go",
	"nft/nft_types.gen.go",
	"orderbook/orderbook_types.gen.go",
	"portfolio/portfolio_types.gen.go",
	"spotprices/spotprices_types.gen.go",
	"tokens/tokens_types.gen.go",
	"traces/traces_types.gen.go",
	"txbroadcast/txbroadcast_types.gen.go",
	"web3/web3_types.gen.go",
}

// requireTools skips the test when a required external tool is unavailable.
func requireTools(t *testing.T) {
	t.Helper()
	for _, tool := range []string{"bash", "jq", "sed", "go"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("skipping: %s not found on PATH", tool)
		}
	}
}

// pipelineEnv returns the test process environment with the oapi-codegen
// install location appended to PATH. The script installs the generator with
// `go install`, which honors GOBIN when set (the Makefile sets it) and falls
// back to GOPATH/bin.
func pipelineEnv() ([]string, error) {
	var binDirs []string
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		binDirs = append(binDirs, gobin)
	}
	out, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return nil, err
	}
	if gopath := strings.TrimSpace(string(out)); gopath != "" {
		binDirs = append(binDirs, filepath.Join(gopath, "bin"))
	}

	env := os.Environ()
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = kv + string(os.PathListSeparator) + strings.Join(binDirs, string(os.PathListSeparator))
			return env, nil
		}
	}
	return append(env, "PATH="+strings.Join(binDirs, string(os.PathListSeparator))), nil
}

// copyFile copies a single file, creating parent directories as needed.
func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// setupPipelineDir builds a temporary repo layout containing the script,
// mapping.json, an openapi directory with the given spec files, and an empty
// sdk-clients directory. specs maps destination file names to source paths.
func setupPipelineDir(t *testing.T, specs map[string]string) string {
	t.Helper()
	root := t.TempDir()
	codegenDir := filepath.Join(root, "codegen")
	require.NoError(t, os.MkdirAll(filepath.Join(codegenDir, "openapi"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "sdk-clients"), 0o755))
	require.NoError(t, copyFile("generate_types.sh", filepath.Join(codegenDir, "generate_types.sh")))
	require.NoError(t, copyFile("mapping.json", filepath.Join(codegenDir, "mapping.json")))
	for name, src := range specs {
		require.NoError(t, copyFile(src, filepath.Join(codegenDir, "openapi", name)))
	}
	return root
}

// runPipeline executes generate_types.sh inside root/codegen and returns its
// combined output.
func runPipeline(root string) (string, error) {
	env, err := pipelineEnv()
	if err != nil {
		return "", err
	}
	cmd := exec.Command("bash", "generate_types.sh")
	cmd.Dir = filepath.Join(root, "codegen")
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// The reproducibility run over all committed specs is expensive (~5s), so it
// is executed once and shared by the tests that assert on its results.
var (
	reproOnce sync.Once
	reproRoot string
	reproOut  string
	reproErr  error
)

func runReproPipeline(t *testing.T) string {
	t.Helper()
	requireTools(t)
	reproOnce.Do(func() {
		root, err := os.MkdirTemp("", "codegen-repro")
		if err != nil {
			reproErr = err
			return
		}
		reproRoot = root
		codegenDir := filepath.Join(root, "codegen")
		files := map[string]string{
			filepath.Join(codegenDir, "generate_types.sh"): "generate_types.sh",
			filepath.Join(codegenDir, "mapping.json"):      "mapping.json",
		}
		for _, name := range specFiles {
			files[filepath.Join(codegenDir, "openapi", name)] = filepath.Join("openapi", name)
		}
		for dst, src := range files {
			if reproErr = copyFile(src, dst); reproErr != nil {
				return
			}
		}
		if reproErr = os.MkdirAll(filepath.Join(root, "sdk-clients"), 0o755); reproErr != nil {
			return
		}
		reproOut, reproErr = runPipeline(root)
	})
	require.NoError(t, reproErr, "pipeline failed:\n%s", reproOut)
	return reproRoot
}

func TestMain(m *testing.M) {
	code := m.Run()
	if reproRoot != "" {
		os.RemoveAll(reproRoot)
	}
	os.Exit(code)
}

// TestGenerateTypesReproducesCommittedFiles is the primary safety net for
// restructuring the codegen pipeline: running the pipeline on the committed
// specs must reproduce every committed *_types.gen.go file byte-for-byte, and
// must produce no other files.
func TestGenerateTypesReproducesCommittedFiles(t *testing.T) {
	root := runReproPipeline(t)

	type testCase struct {
		name string
		file string
	}
	var tests []testCase
	for _, f := range generatedFiles {
		tests = append(tests, testCase{name: f, file: f})
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			generated, err := os.ReadFile(filepath.Join(root, "sdk-clients", tc.file))
			require.NoError(t, err, "pipeline did not produce %s", tc.file)
			committed, err := os.ReadFile(filepath.Join("..", "sdk-clients", tc.file))
			require.NoError(t, err)
			require.Equal(t, string(committed), string(generated),
				"generated output differs from the committed file; if the pipeline change is intentional, regenerate and commit %s", tc.file)
		})
	}

	t.Run("No unexpected files are generated", func(t *testing.T) {
		expected := make(map[string]bool, len(generatedFiles))
		for _, f := range generatedFiles {
			expected[f] = true
		}
		var extra []string
		clientsDir := filepath.Join(root, "sdk-clients")
		err := filepath.WalkDir(clientsDir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, relErr := filepath.Rel(clientsDir, path)
			if relErr != nil {
				return relErr
			}
			if !expected[filepath.ToSlash(rel)] {
				extra = append(extra, rel)
			}
			return nil
		})
		require.NoError(t, err)
		assert.Empty(t, extra, "pipeline produced files not listed in generatedFiles")
	})
}

// TestGenerateTypesSpecsAreFixedPoint asserts that the committed spec files are
// a fixed point of the in-place jq/sed transforms: re-running the pipeline must
// not modify them. This guards against non-idempotent transforms and against
// environment-dependent jq behavior silently rewriting (or deleting parts of)
// the specs, which has happened before with jq 1.6's `|= f?` key deletion.
func TestGenerateTypesSpecsAreFixedPoint(t *testing.T) {
	root := runReproPipeline(t)

	type testCase struct {
		name string
		file string
	}
	var tests []testCase
	for _, f := range specFiles {
		tests = append(tests, testCase{name: f, file: f})
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			after, err := os.ReadFile(filepath.Join(root, "codegen", "openapi", tc.file))
			require.NoError(t, err)
			committed, err := os.ReadFile(filepath.Join("openapi", tc.file))
			require.NoError(t, err)
			require.Equal(t, string(committed), string(after),
				"the pipeline mutated a committed spec; transforms must be idempotent on already-transformed specs")
		})
	}
}

// TestGenerateTypesTransforms documents the observable behavior of each spec
// transform and post-processing step through synthetic fixture specs in
// testdata. Positive assertions are regular expressions because gofmt's field
// alignment depends on neighboring fields.
func TestGenerateTypesTransforms(t *testing.T) {
	requireTools(t)
	root := setupPipelineDir(t, map[string]string{
		"transforms-openapi.json":       filepath.Join("testdata", "transforms-openapi.json"),
		"transforms_extra-openapi.json": filepath.Join("testdata", "transforms_extra-openapi.json"),
	})
	out, err := runPipeline(root)
	require.NoError(t, err, "pipeline failed:\n%s", out)

	mainBytes, err := os.ReadFile(filepath.Join(root, "sdk-clients", "transforms", "transforms_types.gen.go"))
	require.NoError(t, err)
	extraBytes, err := os.ReadFile(filepath.Join(root, "sdk-clients", "transforms", "transforms_extra_types.gen.go"))
	require.NoError(t, err)
	mainSrc := string(mainBytes)
	extraSrc := string(extraBytes)

	tests := []struct {
		name        string
		src         string
		pattern     string
		notContains string
	}{
		{
			name:    "update_operation_ids renames params types via mapping.json",
			src:     mainSrc,
			pattern: `type GetQuoteParams struct`,
		},
		{
			name:        "update_operation_ids removes the raw controller name",
			src:         mainSrc,
			notContains: `AggregationControllerGetQuoteParams`,
		},
		{
			name:    "fix_chain_id_number_types generates chain-id parameters as int",
			src:     mainSrc,
			pattern: "ChainId int `url:\"chainId\" json:\"chainId\"`",
		},
		{
			name:    "fix_chain_id_number_types generates chain-id properties as int",
			src:     mainSrc,
			pattern: "SrcChainId int\\s+`json:\"srcChainId\"`",
		},
		{
			name:    "non-chain number fields remain float32",
			src:     mainSrc,
			pattern: "Price float32\\s+`json:\"price\"`",
		},
		{
			name:    "check_and_fix_incorrect_number_arrays rewrites number\\[\\] to a float32 slice",
			src:     mainSrc,
			pattern: "Amounts \\[\\]float32 `url:\"amounts\" json:\"amounts\"`",
		},
		{
			name:    "change_any_of_int_to_int collapses anyOf parameters to the first variant",
			src:     mainSrc,
			pattern: "GasPrice int `url:\"gasPrice,omitempty\" json:\"gasPrice,omitempty\"`",
		},
		{
			name:    "remove_null_schemas drops null anyOf variants before the collapse",
			src:     mainSrc,
			pattern: "Filter string `url:\"filter,omitempty\" json:\"filter,omitempty\"`",
		},
		{
			name:    "add_pointer_skip_field keeps optional parameters pointer-free",
			src:     mainSrc,
			pattern: "Note string `url:\"note,omitempty\" json:\"note,omitempty\"`",
		},
		{
			name:    "add_pointer_skip_field keeps optional typed properties pointer-free",
			src:     mainSrc,
			pattern: "OptionalNote string `json:\"optionalNote,omitempty\"`",
		},
		{
			name:    "change_any_of_ref_to_ref collapses allOf-wrapped refs; bare refs keep the optional pointer",
			src:     mainSrc,
			pattern: "Linked \\*Linked\\s+`json:\"linked,omitempty\"`",
		},
		{
			name:        "form tags are rewritten to url tags",
			src:         mainSrc,
			notContains: "`form:\"",
		},
		{
			name:    "underscored spec names generate into the base package",
			src:     extraSrc,
			pattern: `package transforms\n`,
		},
		{
			name:    "underscored spec names keep the full name in the generated file",
			src:     extraSrc,
			pattern: "Id string `json:\"id\"`",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.pattern != "" {
				assert.Regexp(t, regexp.MustCompile(tc.pattern), tc.src)
			}
			if tc.notContains != "" {
				assert.NotContains(t, tc.src, tc.notContains)
			}
		})
	}
}

// TestGenerateTypesRejectsInvalidSpecFilenames pins the naming contract: every
// file in the openapi directory must end in -openapi.json or the script aborts.
func TestGenerateTypesRejectsInvalidSpecFilenames(t *testing.T) {
	requireTools(t)

	tests := []struct {
		name     string
		specName string
		errorMsg string
	}{
		{
			name:     "Missing -openapi suffix",
			specName: "transforms.json",
			errorMsg: "does not match the expected naming schema",
		},
		{
			name:     "Wrong extension",
			specName: "transforms-openapi.yaml",
			errorMsg: "does not match the expected naming schema",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := setupPipelineDir(t, map[string]string{
				tc.specName: filepath.Join("testdata", "transforms-openapi.json"),
			})
			out, err := runPipeline(root)
			require.Error(t, err, "script should reject %s:\n%s", tc.specName, out)
			assert.Contains(t, out, tc.errorMsg)
		})
	}
}
