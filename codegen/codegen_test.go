// Characterization tests for the type-generation pipeline. They pin down the
// pipeline's observable behavior so it can be safely evolved: any change must
// keep these tests green, or change them in a reviewed commit alongside the
// regenerated files.
//
// The pipeline invokes the pinned oapi-codegen via `go run module@version`,
// which needs the module cache or network on first run.
package codegen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
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

// The reproducibility run over all committed specs is expensive, so it is
// executed once and shared by the tests that assert on its results.
var (
	reproOnce sync.Once
	reproOut  string
	reproErr  error
)

func runReproPipeline(t *testing.T) string {
	t.Helper()
	reproOnce.Do(func() {
		dir, err := os.MkdirTemp("", "codegen-repro")
		if err != nil {
			reproErr = err
			return
		}
		reproOut = dir
		reproErr = Generate(Options{
			SpecsDir:    "openapi",
			MappingFile: "mapping.json",
			OutputDir:   dir,
		})
	})
	require.NoError(t, reproErr, "pipeline failed")
	return reproOut
}

func TestMain(m *testing.M) {
	code := m.Run()
	if reproOut != "" {
		os.RemoveAll(reproOut)
	}
	os.Exit(code)
}

// TestGenerateTypesReproducesCommittedFiles is the primary safety net for
// evolving the codegen pipeline: running the pipeline on the committed specs
// must reproduce every committed *_types.gen.go file byte-for-byte, produce no
// other files, and leave the committed spec files untouched.
func TestGenerateTypesReproducesCommittedFiles(t *testing.T) {
	specsBefore := readAll(t, "openapi", specFiles)
	out := runReproPipeline(t)

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
			generated, err := os.ReadFile(filepath.Join(out, tc.file))
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
		err := filepath.WalkDir(out, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, relErr := filepath.Rel(out, path)
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

	t.Run("Committed specs are not modified", func(t *testing.T) {
		specsAfter := readAll(t, "openapi", specFiles)
		require.Equal(t, specsBefore, specsAfter, "Generate must never modify the spec files on disk")
	})
}

// TestTransformsAreIdempotent asserts that applying the transform chain to an
// already-transformed spec changes nothing. This guards against transforms
// that corrupt their own output — the committed specs have all historically
// been normalized in place, so every transform must be a no-op on them.
func TestTransformsAreIdempotent(t *testing.T) {
	mapping, err := loadOperationIDMapping("mapping.json")
	require.NoError(t, err)

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
			raw, err := os.ReadFile(filepath.Join("openapi", tc.file))
			require.NoError(t, err)
			once, err := applyTransforms(tc.file, raw, mapping)
			require.NoError(t, err)
			twice, err := applyTransforms(tc.file, once, mapping)
			require.NoError(t, err)

			var first, second any
			require.NoError(t, json.Unmarshal(once, &first))
			require.NoError(t, json.Unmarshal(twice, &second))
			require.Equal(t, first, second, "transform chain is not idempotent")
		})
	}
}

// TestGenerateTypesTransforms documents the observable behavior of each spec
// transform and post-processing step through synthetic fixture specs in
// testdata. Positive assertions are regular expressions because gofmt's field
// alignment depends on neighboring fields.
func TestGenerateTypesTransforms(t *testing.T) {
	out := t.TempDir()
	err := Generate(Options{
		SpecsDir:    "testdata",
		MappingFile: "mapping.json",
		OutputDir:   out,
	})
	require.NoError(t, err, "pipeline failed on fixture specs")

	mainBytes, err := os.ReadFile(filepath.Join(out, "transforms", "transforms_types.gen.go"))
	require.NoError(t, err)
	extraBytes, err := os.ReadFile(filepath.Join(out, "transforms", "transforms_extra_types.gen.go"))
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
			name:    "renameOperationIDs renames params types via mapping.json",
			src:     mainSrc,
			pattern: `type GetQuoteParams struct`,
		},
		{
			name:        "renameOperationIDs removes the raw controller name",
			src:         mainSrc,
			notContains: `AggregationControllerGetQuoteParams`,
		},
		{
			name:    "fixChainIDNumberTypes generates chain-id parameters as int",
			src:     mainSrc,
			pattern: "ChainId int `url:\"chainId\" json:\"chainId\"`",
		},
		{
			name:    "fixChainIDNumberTypes generates chain-id properties as int",
			src:     mainSrc,
			pattern: "SrcChainId int\\s+`json:\"srcChainId\"`",
		},
		{
			name:    "non-chain number fields remain float32",
			src:     mainSrc,
			pattern: "Price float32\\s+`json:\"price\"`",
		},
		{
			name:    "fixIncorrectNumberArrays rewrites number\\[\\] to a float32 slice",
			src:     mainSrc,
			pattern: "Amounts \\[\\]float32 `url:\"amounts\" json:\"amounts\"`",
		},
		{
			name:    "collapseAnyOfParameters collapses anyOf parameters to the first variant",
			src:     mainSrc,
			pattern: "GasPrice int `url:\"gasPrice,omitempty\" json:\"gasPrice,omitempty\"`",
		},
		{
			name:    "removeNullAnyOfVariants drops null anyOf variants before the collapse",
			src:     mainSrc,
			pattern: "Filter string `url:\"filter,omitempty\" json:\"filter,omitempty\"`",
		},
		{
			name:    "addSkipOptionalPointer keeps optional parameters pointer-free",
			src:     mainSrc,
			pattern: "Note string `url:\"note,omitempty\" json:\"note,omitempty\"`",
		},
		{
			name:    "addSkipOptionalPointer keeps optional typed properties pointer-free",
			src:     mainSrc,
			pattern: "OptionalNote string `json:\"optionalNote,omitempty\"`",
		},
		{
			name:    "simplifyAllOfRefs collapses allOf-wrapped refs; bare refs keep the optional pointer",
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
// file in the specs directory must end in -openapi.json or generation aborts.
func TestGenerateTypesRejectsInvalidSpecFilenames(t *testing.T) {
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
			specsDir := t.TempDir()
			fixture, err := os.ReadFile(filepath.Join("testdata", "transforms-openapi.json"))
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(filepath.Join(specsDir, tc.specName), fixture, 0o644))

			err = Generate(Options{
				SpecsDir:    specsDir,
				MappingFile: "mapping.json",
				OutputDir:   t.TempDir(),
			})
			require.Error(t, err, "generation should reject %s", tc.specName)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

// TestSpecSourcesCoverAllSpecs pins the fetch manifest to the committed spec
// set: every spec must have a manifest entry (possibly empty, meaning "source
// not confirmed yet"), and the manifest must not reference unknown specs.
func TestSpecSourcesCoverAllSpecs(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{}
	for _, f := range specFiles {
		tests = append(tests, struct {
			name string
			file string
		}{name: f, file: f})
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := specSources[tc.file]
			assert.True(t, ok, "spec %s has no entry in specSources (codegen/fetch.go)", tc.file)
		})
	}

	t.Run("No unknown manifest entries", func(t *testing.T) {
		known := make(map[string]bool, len(specFiles))
		for _, f := range specFiles {
			known[f] = true
		}
		for name := range specSources {
			assert.True(t, known[name], "specSources entry %s does not match any committed spec", name)
		}
	})
}

// TestSpecsMatchLock verifies the provenance chain: every committed spec file
// must hash to exactly what specs.lock.json records. Together with
// TestGenerateTypesReproducesCommittedFiles this ties the generated code to a
// specific upstream snapshot: generated code ⇔ committed specs ⇔ lock entry
// (source, upstream version, hash, fetch time).
//
// If this test fails, either a spec was hand-edited (not allowed — put
// corrections in codegen/overrides.go) or specs were updated without the fetch
// tool (run `go run ./codegen/cmd/fetch-specs`, or `-seed` to repair).
func TestSpecsMatchLock(t *testing.T) {
	lock, err := loadSpecLock("specs.lock.json")
	require.NoError(t, err)

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
			entry, ok := lock[tc.file]
			require.True(t, ok, "spec %s has no entry in specs.lock.json", tc.file)
			raw, err := os.ReadFile(filepath.Join("openapi", tc.file))
			require.NoError(t, err)
			assert.Equal(t, entry.SHA256, sha256Hex(raw),
				"committed spec does not match specs.lock.json — hand edits are not allowed; use codegen/overrides.go for corrections and the fetch-specs tool for updates")
			assert.Equal(t, entry.Version, specVersion(raw),
				"lock version does not match the spec's info.version")
		})
	}

	t.Run("No unknown lock entries", func(t *testing.T) {
		known := make(map[string]bool, len(specFiles))
		for _, f := range specFiles {
			known[f] = true
		}
		for name := range lock {
			assert.True(t, known[name], "lock entry %s does not match any committed spec", name)
		}
	})
}

// TestValidateAndIndentSpec documents the fetch-time validation and
// normalization contract: whitespace is normalized, key order and values are
// preserved, and non-OpenAPI payloads are rejected.
func TestValidateAndIndentSpec(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:     "Valid spec is re-indented preserving key order",
			input:    `{"openapi":"3.0.0","zebra":1,"alpha":2,"paths":{}}`,
			expected: "{\n  \"openapi\": \"3.0.0\",\n  \"zebra\": 1,\n  \"alpha\": 2,\n  \"paths\": {}\n}\n",
		},
		{
			name:        "Missing openapi field",
			input:       `{"paths":{}}`,
			expectError: true,
			errorMsg:    `missing "openapi"`,
		},
		{
			name:        "Missing paths object",
			input:       `{"openapi":"3.0.0"}`,
			expectError: true,
			errorMsg:    `missing "paths"`,
		},
		{
			name:        "HTML error page is rejected",
			input:       `<!doctype html><html></html>`,
			expectError: true,
			errorMsg:    "invalid spec JSON",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := validateAndIndentSpec([]byte(tc.input))
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, string(out))
			}
		})
	}
}

func readAll(t *testing.T, dir string, names []string) map[string]string {
	t.Helper()
	files := make(map[string]string, len(names))
	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		require.NoError(t, err)
		files[name] = string(raw)
	}
	return files
}
