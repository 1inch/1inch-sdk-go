package http_executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"regexp"

	"github.com/google/go-querystring/query"

	"github.com/1inch/1inch-sdk-go/common"
)

func DefaultHttpClient(apiUrl string, apiKey string) (*Client, error) {
	baseURL, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}
	return &Client{
		httpClient: *http.DefaultClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}, nil
}

type Client struct {
	httpClient http.Client

	// The URL of the 1inch API
	baseURL *url.URL
	// The API key to use for authentication
	apiKey string
}

func (c *Client) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v interface{}) error {
	u, err := addQueryParameters(payload.U, payload.Params)
	if err != nil {
		return err
	}

	fullURL, err := c.baseURL.Parse(u)
	if err != nil {
		return err
	}

	req, err := c.prepareRequest(ctx, payload.Method, fullURL, payload.Body)
	if err != nil {
		return fmt.Errorf("preparing request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request failed: %w", err)
	}

	if err = c.processResponse(resp, v); err != nil {
		return fmt.Errorf("processing response failed: %w", err)
	}

	return nil
}

func (c *Client) prepareRequest(ctx context.Context, method string, fullURL *url.URL, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "1inch-dev-portal-client-go:beta.2")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) processResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return c.handleErrorResponse(resp)
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	if buf.Len() == 0 {
		return nil // No content to decode
	}

	if v != nil {
		return json.NewDecoder(buf).Decode(v)
	}

	return nil
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading error response body failed: %v", err)
	}

	var errorMessageMap map[string]interface{}
	err = json.Unmarshal(body, &errorMessageMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal error response body: %v", err)
	}

	errFormatted, err := json.MarshalIndent(errorMessageMap, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to format error response body: %v", err)
	}

	return fmt.Errorf("%s", errFormatted)
}

// addQueryParameters adds the parameters in the struct params as URL query parameters to s.
// params must be a struct whose fields may contain "url" tags.
func addQueryParameters(s string, params interface{}) (string, error) {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(params)
	if err != nil {
		return s, err
	}

	for k, v := range qs {
		if isScientificNotation(v[0]) {
			expanded, err := expandScientificNotation(v[0])
			if err != nil {
				return "", fmt.Errorf("failed to expand scientific notation for parameter %v with a value of %v: %v", k, v, err)
			}
			v[0] = expanded
		}
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// isScientificNotation checks if the string is in scientific notation (like 1e+18).
func isScientificNotation(s string) bool {
	// This regular expression matches strings in the format of "1e+18", "2.3e-4", etc.
	re := regexp.MustCompile(`^[+-]?\d+(\.\d+)?[eE][+-]?\d+$`)
	return re.MatchString(s)
}

func expandScientificNotation(s string) (string, error) {
	f, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
	if err != nil {
		return "", err
	}

	// Use a precision that is sufficient to handle small numbers.
	// The precision here is set to a large number to ensure accuracy for small decimal values.
	f.SetPrec(64)

	return f.Text('f', -1), nil // -1 ensures that insignificant zeroes are not omitted
}
