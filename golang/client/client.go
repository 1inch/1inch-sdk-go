package client

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
)

// This is the base URL for the 1inch API.
var baseUrlProduction, _ = url.Parse("http://api.1inch.dev")
var baseUrlStaging, _ = url.Parse("http://fake-staging.1inch.dev")

type Environment string

const (
	EnvironmentProduction Environment = "Production"
	EnvironmentStaging    Environment = "Staging"
)

type service struct {
	client *Client
}

type Config struct {
	TargetEnvironment Environment
	ApiKey            string
}

func NewClient(config Config) (*Client, error) {

	var baseUrl *url.URL
	switch config.TargetEnvironment {
	case "":
		fallthrough
	case EnvironmentProduction:
		baseUrl = baseUrlProduction
	case EnvironmentStaging:
		baseUrl = baseUrlStaging
	}

	c := &Client{
		httpClient: &http.Client{},
		BaseURL:    baseUrl,
		ApiKey:     config.ApiKey,
	}

	c.common.client = c

	c.Swap = (*SwapService)(&c.common)
	c.TokenPrices = (*TokenPricesService)(&c.common)

	return c, nil
}

type Client struct {
	httpClient *http.Client

	BaseURL *url.URL
	ApiKey  string

	common service

	Swap        *SwapService
	TokenPrices *TokenPricesService
}

type ErrorResponse struct {
	Response     *http.Response `json:"-"`
	ErrorMessage string         `json:"error"`
	Description  string         `json:"description"`
	StatusCode   int            `json:"statusCode"`
	Meta         []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"meta"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v - %v - %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.ErrorMessage, r.Description, r.Meta)
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))

	req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	// Check response codes
	var errorResp *ErrorResponse
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorResp = &ErrorResponse{Response: resp}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, errorResp)
		if err != nil {
			// reset the response as if this never happened
			errorResp = &ErrorResponse{Response: resp}
		}
	}
	if errorResp != nil {
		return nil, errorResp
	}

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = fmt.Errorf("request did not fail, but the response could not be decoded (this could be due to cloudflare blocking the request): %v", decErr)
		}
	}
	return resp, err
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// TODO verify if this is needed
	//if !strings.HasSuffix(c.BaseURL.Path, "/") {
	//	return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	//}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
