package codegen

import (
	"bytes"
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
// An empty URL means the source has not been confirmed yet: the fusion and
// fusion-plus services are split into orders/quoter/relayer sub-APIs whose
// portal slugs cannot be discovered without a token (the endpoint returns 401
// before routing). Fill these in from the portal documentation page's network
// requests, then delete this note once all sources are confirmed.
var specSources = map[string]string{
	"aggregation-openapi.json":        "https://api.1inch.dev/portal/apis/swagger/swap",
	"balances-openapi.json":           "https://api.1inch.dev/portal/apis/swagger/balance",
	"fusion_orders-openapi.json":      "",
	"fusion_quoter-openapi.json":      "",
	"fusion_relayer-openapi.json":     "",
	"fusionplus_orders-openapi.json":  "",
	"fusionplus_quoter-openapi.json":  "",
	"fusionplus_relayer-openapi.json": "",
	"gasprices-openapi.json":          "https://api.1inch.dev/portal/apis/swagger/gas-price",
	"history-openapi.json":            "https://api.1inch.dev/portal/apis/swagger/history",
	"nft-openapi.json":                "https://api.1inch.dev/portal/apis/swagger/nft",
	"orderbook-openapi.json":          "https://api.1inch.dev/portal/apis/swagger/orderbook",
	"portfolio-openapi.json":          "https://api.1inch.dev/portal/apis/swagger/portfolio",
	"spotprices-openapi.json":         "https://api.1inch.dev/portal/apis/swagger/spot-price",
	"tokens-openapi.json":             "https://api.1inch.dev/portal/apis/swagger/token",
	"traces-openapi.json":             "https://api.1inch.dev/portal/apis/swagger/traces",
	"txbroadcast-openapi.json":        "https://api.1inch.dev/portal/apis/swagger/tx-gateway",
	"web3-openapi.json":               "https://api.1inch.dev/portal/apis/swagger/web3",
}

// FetchOptions configures a spec refresh.
type FetchOptions struct {
	// SpecsDir is the directory the fetched specs are written to.
	SpecsDir string
	// MappingFile is the operation-id mapping file, used to dry-run the
	// transform pipeline against each fetched spec.
	MappingFile string
	// Token is the Dev Portal API token.
	Token string
	// Client is the HTTP client to use; defaults to a client with a 30s
	// timeout.
	Client *http.Client
}

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

	var errs []error
	for _, name := range names {
		url := specSources[name]
		if url == "" {
			fmt.Printf("skipped   %s (no source configured in codegen/fetch.go)\n", name)
			continue
		}
		if err := fetchOne(client, opts, name, url, mapping); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			fmt.Printf("FAILED    %s: %v\n", name, err)
			continue
		}
		fmt.Printf("refreshed %s\n", name)
	}
	return errors.Join(errs...)
}

func fetchOne(client *http.Client, opts FetchOptions, name, url string, mapping map[string]string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+opts.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s: %.200s", resp.StatusCode, url, body)
	}

	normalized, err := validateAndIndentSpec(body)
	if err != nil {
		return fmt.Errorf("response from %s is not a usable OpenAPI document: %w", url, err)
	}

	// Dry-run the transform pipeline so schema drift that breaks an override
	// is caught at fetch time, not at the next generation.
	if _, err := applyTransforms(name, normalized, mapping); err != nil {
		return fmt.Errorf("fetched spec no longer matches the codegen overrides (update codegen/overrides.go): %w", err)
	}

	return os.WriteFile(filepath.Join(opts.SpecsDir, name), normalized, 0o644)
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
