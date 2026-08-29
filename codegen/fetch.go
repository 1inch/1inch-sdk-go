package codegen

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// specSources maps each committed spec file to the Dev Portal endpoint that
// serves it. The endpoints require a Dev Portal token (Authorization: Bearer).
//
// The committed spec files are meant to be untouched copies of these upstream
// documents; all corrections happen at generation time. Refresh workflow:
//
//	DEV_PORTAL_TOKEN=... make codegen-fetch-specs   # rewrite local copies
//	git diff codegen/openapi                        # review upstream changes
//	make codegen-types                              # regenerate (overrides fail loudly on mismatches)
//	go test ./codegen                               # re-pin the output
//
// An empty URL means the API has no fetchable upstream source right now; see
// the per-entry comments.
var specSources = map[string]string{
	// Confirmed against the live portal on 2026-08-28. Sources whose URL embeds
	// a version (history, traces, portfolio) pin that version and need a bump
	// here when upstream releases a new major.
	"aggregation-openapi.json":        "https://api.1inch.dev/swap/swagger/ethereum-json",
	"balances-openapi.json":           "https://api.1inch.dev/balance/swagger/ethereum-json",
	"fusion_orders-openapi.json":      "https://api.1inch.dev/fusion/orders/swagger/ethereum-json",
	"fusion_quoter-openapi.json":      "https://api.1inch.dev/fusion/quoter/swagger/ethereum-json",
	"fusion_relayer-openapi.json":     "https://api.1inch.dev/fusion/relayer/swagger/ethereum-json",
	"fusionplus_orders-openapi.json":  "https://api.1inch.dev/fusion-plus/orders/swagger-json",
	"fusionplus_quoter-openapi.json":  "https://api.1inch.dev/fusion-plus/quoter/swagger-json",
	"fusionplus_relayer-openapi.json": "https://api.1inch.dev/fusion-plus/relayer/swagger-json",
	"gasprices-openapi.json":          "https://api.1inch.dev/gas-price/swagger/ethereum-json",
	"history-openapi.json":            "https://api.1inch.dev/history/v2.0/history/1/swagger-json",
	"nft-openapi.json":                "https://api.1inch.dev/nft/swagger-json",
	// The portal documentation references /orderbook/swagger/ethereum-json for
	// this API, but that endpoint currently returns 404; left unconfigured
	// until upstream fixes it.
	"orderbook-openapi.json":   "",
	"portfolio-openapi.json":   "https://api.1inch.dev/portfolio/v5.0/openapi.json",
	"spotprices-openapi.json":  "https://api.1inch.dev/price/swagger/ethereum-json",
	"tokens-openapi.json":      "https://api.1inch.dev/token/swagger-json",
	"traces-openapi.json":      "https://api.1inch.dev/traces/v1.0/chain/1/swagger-json",
	"txbroadcast-openapi.json": "https://api.1inch.dev/tx-gateway/swagger/ethereum-json",
	// The web3 API has no spec in the portal documentation catalog; the
	// committed spec appears to be authored for this SDK and has no upstream
	// source to fetch.
	"web3-openapi.json": "",
}

// FetchOptions configures a spec refresh.
type FetchOptions struct {
	// SpecsDir is the directory the fetched specs are written to.
	SpecsDir string
	// MappingFile is the operation-id mapping file, used to dry-run the
	// transform pipeline against each fetched spec.
	MappingFile string
	// LockFile records the provenance of each committed spec; see specLock.
	LockFile string
	// Token is the Dev Portal API token.
	Token string
	// Client is the HTTP client to use; defaults to a client with a 30s
	// timeout.
	Client *http.Client
	// Now stamps fetchedAt in the lock file; defaults to time.Now.
	Now func() time.Time
}

// specLockEntry ties one committed spec file to the upstream snapshot it was
// copied from. The lock file (specs.lock.json) is the provenance record for
// the whole chain: generated code is pinned to the committed specs by the
// characterization tests, and the committed specs are pinned to upstream by
// this file. CI verifies the hashes, so hand-edited spec copies fail the
// build.
type specLockEntry struct {
	// Source is the URL the spec was fetched from.
	Source string `json:"source,omitempty"`
	// Version is the upstream document's info.version field.
	Version string `json:"version,omitempty"`
	// SHA256 is the hex digest of the committed spec file.
	SHA256 string `json:"sha256"`
	// FetchedAt is when the spec was last fetched (RFC 3339). Empty for
	// entries seeded from pre-existing local copies whose fetch time is
	// unknown.
	FetchedAt string `json:"fetchedAt,omitempty"`
}

type specLock map[string]specLockEntry

// FetchSpecs downloads every configured spec from the Dev Portal and rewrites
// the local copies. Each fetched document is validated (it must be an OpenAPI
// document and must survive the transform pipeline, so broken overrides
// surface immediately) and re-indented with json.Indent, which changes
// whitespace only — key order and values stay byte-exact as served.
//
// Specs without a configured source are reported and skipped. All configured
// specs are attempted; errors are aggregated so one failure does not hide the
// rest.
func FetchSpecs(opts FetchOptions) error {
	if opts.Token == "" {
		return errors.New("a Dev Portal token is required (set DEV_PORTAL_TOKEN or pass -token)")
	}
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	mapping, err := loadOperationIDMapping(opts.MappingFile)
	if err != nil {
		return err
	}

	names := make([]string, 0, len(specSources))
	for name := range specSources {
		names = append(names, name)
	}
	sort.Strings(names)

	lock, err := loadSpecLock(opts.LockFile)
	if err != nil {
		return err
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}

	var errs []error
	for _, name := range names {
		url := specSources[name]
		if url == "" {
			fmt.Printf("skipped   %s (no source configured in codegen/fetch.go)\n", name)
			continue
		}
		normalized, err := fetchOne(client, opts, name, url, mapping)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			fmt.Printf("FAILED    %s: %v\n", name, err)
			continue
		}
		entry := specLockEntry{
			Source:    url,
			Version:   specVersion(normalized),
			SHA256:    sha256Hex(normalized),
			FetchedAt: now().UTC().Format(time.RFC3339),
		}
		if prev, ok := lock[name]; ok && prev.SHA256 == entry.SHA256 {
			// Preserve the prior fetch time when the content is unchanged, so the
			// lock only changes on a real content/version/source change. Otherwise
			// every fetch would rewrite timestamps and the weekly spec-drift
			// workflow would open a false-positive PR containing only timestamp
			// noise.
			entry.FetchedAt = prev.FetchedAt
			fmt.Printf("unchanged %s (version %s)\n", name, entry.Version)
		} else {
			fmt.Printf("refreshed %s (version %s)\n", name, entry.Version)
		}
		lock[name] = entry
	}

	if err := writeSpecLock(opts.LockFile, lock); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// SeedLock builds the lock file from the spec copies already on disk. It is
// the bootstrap (and repair) path for specs that predate the fetch tooling:
// provenance fields that cannot be known locally (fetchedAt) are left empty.
func SeedLock(opts FetchOptions) error {
	lock := specLock{}
	names := make([]string, 0, len(specSources))
	for name := range specSources {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(opts.SpecsDir, name))
		if err != nil {
			return err
		}
		lock[name] = specLockEntry{
			Source:  specSources[name],
			Version: specVersion(raw),
			SHA256:  sha256Hex(raw),
		}
		fmt.Printf("seeded %s (version %s)\n", name, lock[name].Version)
	}
	return writeSpecLock(opts.LockFile, lock)
}

func loadSpecLock(path string) (specLock, error) {
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return specLock{}, nil
	}
	if err != nil {
		return nil, err
	}
	var lock specLock
	if err := json.Unmarshal(raw, &lock); err != nil {
		return nil, fmt.Errorf("invalid lock file %s: %w", path, err)
	}
	return lock, nil
}

func writeSpecLock(path string, lock specLock) error {
	raw, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}

// specVersion extracts the upstream document's info.version, or "" if absent.
func specVersion(raw []byte) string {
	doc, err := decodeSpec(raw)
	if err != nil {
		return ""
	}
	info, _ := doc["info"].(spec)
	version, _ := info["version"].(string)
	return version
}

func sha256Hex(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func fetchOne(client *http.Client, opts FetchOptions, name, url string, mapping map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+opts.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s: %.200s", resp.StatusCode, url, body)
	}

	normalized, err := validateAndIndentSpec(body)
	if err != nil {
		return nil, fmt.Errorf("response from %s is not a usable OpenAPI document: %w", url, err)
	}

	// Dry-run the transform pipeline so schema drift that breaks an override
	// is caught at fetch time, not at the next generation.
	if _, err := applyTransforms(name, normalized, mapping); err != nil {
		return nil, fmt.Errorf("fetched spec no longer matches the codegen overrides (update codegen/overrides.go): %w", err)
	}

	if err := os.WriteFile(filepath.Join(opts.SpecsDir, name), normalized, 0o644); err != nil {
		return nil, err
	}
	return normalized, nil
}

// validateAndIndentSpec checks that the payload is an OpenAPI document and
// normalizes its indentation without touching key order or values.
func validateAndIndentSpec(raw []byte) ([]byte, error) {
	doc, err := decodeSpec(raw)
	if err != nil {
		return nil, err
	}
	if _, ok := doc["openapi"].(string); !ok {
		return nil, errors.New(`missing "openapi" version field`)
	}
	if _, ok := doc["paths"].(spec); !ok {
		return nil, errors.New(`missing "paths" object`)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, bytes.TrimSpace(raw), "", "  "); err != nil {
		return nil, err
	}
	out.WriteByte('\n')
	return out.Bytes(), nil
}
