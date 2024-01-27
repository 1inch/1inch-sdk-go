package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/go-querystring/query"

	"github.com/1inch/1inch-sdk/golang/helpers"
)

type Environment string

const (
	EnvironmentProduction Environment = "Production"
	EnvironmentStaging    Environment = "Staging"
)

type service struct {
	client *Client
}

type Config struct {
	DevPortalApiKey   string
	Web3HttpProviders []Web3ProviderConfig
	TenderlyKey       string
}

type Web3ProviderConfig struct {
	ChainId int
	Url     string
}

func (c *Config) validate() error {

	if c.DevPortalApiKey == "" {
		return fmt.Errorf("API key is required")
	}
	if len(c.Web3HttpProviders) == 0 {
		return fmt.Errorf("at least one web3 provider URL is required")
	}

	return nil
}

type Client struct {
	// Standard http client in Go
	httpClient *http.Client
	// Ethereum client map
	EthClientMap map[int]*ethclient.Client
	// The URL of the 1inch API
	ApiBaseURL *url.URL
	// The API key to use for authentication
	ApiKey string
	// When present, tests will simulate swaps on Tenderly
	TenderlyKey string
	// Once a transaction has been sent by the SDK, the nonce is tracked internally to avoid RPC desync issues on subsequent transactions
	NonceCache map[string]uint64
	// A struct that will contain a reference to this client. Used to separate each API into a unique namespace to aid in method discovery
	common service
	// Isolated namespaces for each API
	Actions   *ActionService
	Swap      *SwapService
	Orderbook *OrderbookService
}

func (c *Client) GetEthClient(chainId int) (*ethclient.Client, error) {
	ethClient, ok := c.EthClientMap[chainId]
	if !ok {
		return nil, fmt.Errorf("no client for chain id %d", chainId)
	}
	return ethClient, nil
}

// NewClient creates and initializes a new Client instance based on the provided Config.
func NewClient(config Config) (*Client, error) {
	// TODO this may be replaceable with https://github.com/go-playground/validator
	err := config.validate()
	if err != nil {
		return nil, fmt.Errorf("config validation error: %v", err)
	}

	var ethClient *ethclient.Client
	ethClientMap := make(map[int]*ethclient.Client)
	for _, provider := range config.Web3HttpProviders {
		ethClient, err = ethclient.Dial(provider.Url)
		if err != nil {
			return nil, fmt.Errorf("failed to create eth client: %v", err)
		}
		ethClientMap[provider.ChainId] = ethClient
	}

	apiBaseUrl, err := url.Parse("https://api.1inch.dev")
	if err != nil {
		return nil, fmt.Errorf("failed to parse API base URL: %v", err)
	}

	c := &Client{
		httpClient:   &http.Client{},
		EthClientMap: ethClientMap,
		ApiBaseURL:   apiBaseUrl,
		ApiKey:       config.DevPortalApiKey,
		NonceCache:   make(map[string]uint64),
		TenderlyKey:  config.TenderlyKey,
	}

	c.common.client = c

	c.Actions = (*ActionService)(&c.common)
	c.Swap = (*SwapService)(&c.common)
	c.Orderbook = (*OrderbookService)(&c.common)

	return c, nil
}

func (c *Client) NewRequest(method, urlStr string, body []byte) (*http.Request, error) {
	u, err := c.ApiBaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
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
	// TODO errors are handled generically at the moment
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		// Unmarshal into a map to handle arbitrary JSON structure
		var messageMap map[string]interface{}
		err = json.Unmarshal(data, &messageMap)
		if err != nil {
			// Fallback to raw string if unmarshalling fails
			return nil, fmt.Errorf("failed to unmarshal response body: %s", string(data))
		}

		// Marshal the message with indentation
		formattedMessage, err := json.MarshalIndent(messageMap, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal formatted message: %v - Original error: %s", err, string(data))
		}

		return nil, fmt.Errorf("%s", formattedMessage)
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

// ReplacePathVariable replaces the path variable in the given URL with the specified value.
func ReplacePathVariable(path, pathVarName string, value interface{}) (string, error) {
	placeholder := fmt.Sprintf("{%s}", pathVarName)

	if !strings.Contains(path, placeholder) {
		return "", errors.New("path variable not found in URL path")
	}

	return strings.Replace(path, placeholder, fmt.Sprintf("%s", value), 1), nil
}
