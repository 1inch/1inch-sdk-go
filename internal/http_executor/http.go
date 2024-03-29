package http_executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"

	"github.com/1inch/1inch-sdk-go/internal/common"
	"github.com/1inch/1inch-sdk-go/internal/helpers"
)

func DefaultHttpClient(baseURL *url.URL, apiKey string) Client {
	httpClient := http.Client{}
	return Client{
		httpClient,
		baseURL,
		apiKey,
	}
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

	req.Header.Set("User-Agent", "GolangSDK/0.0.3-developer-preview")
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

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading error response body failed: %v", err)
	}

	var apiError struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiError); err != nil {
		// Fallback to raw body text if the body cannot be parsed as JSON
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, apiError.Message)
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
		if helpers.IsScientificNotation(v[0]) {
			expanded, err := helpers.ExpandScientificNotation(v[0])
			if err != nil {
				return "", fmt.Errorf("failed to expand scientific notation for parameter %v with a value of %v: %v", k, v, err)
			}
			v[0] = expanded
		}
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
