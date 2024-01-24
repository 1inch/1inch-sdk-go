package client

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/go-querystring/query"

	"github.com/1inch/1inch-sdk/golang/helpers"
)

// This is the base URL for the 1inch API.
var baseUrlProduction, _ = url.Parse("https://api.1inch.dev")
var baseUrlStaging, _ = url.Parse("https://fake-staging.1inch.dev")

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
	ChainId           int
	DevPortalApiKey   string
	Web3HttpProvider  string
	WalletKey         string
	TenderlyKey       string
}

func (c *Config) validate() error {

	if c.DevPortalApiKey == "" {
		return fmt.Errorf("API key is required")
	}
	if c.Web3HttpProvider == "" {
		return fmt.Errorf("web3 provider URL is required")
	}
	if c.ChainId == 0 {
		return fmt.Errorf("chain ID is required")
	} else if !helpers.IsValidChainId(c.ChainId) {
		return fmt.Errorf("invalid chain id: %d", c.ChainId)
	}

	return nil
}

type Client struct {
	// Standard http client in Go
	httpClient *http.Client
	// Ethereum client
	EthClient *ethclient.Client
	// The chain ID for requests
	ChainId int
	// The URL of the 1inch API
	BaseURL *url.URL
	// The API key to use for authentication
	ApiKey string
	// The key of the wallet that will be used to sign transactions
	WalletKey string
	// The public address of the wallet that will be used to sign transactions (derived from the private key)
	// DO NOT MANUALLY SET
	PublicAddress common.Address
	// RPC URL for web3 provider with key
	RpcUrlWithKey string
	// Once a transaction has been sent by the SDK, the nonce is tracked internally to avoid RPC desync issues on subsequent transactions
	NonceCache map[string]uint64
	// A struct that will contain a reference to this client. Used to separate each API into a unique namespace to aid in method discovery
	common service
	// Isolated namespaces for each API
	Actions   *ActionService
	Swap      *SwapService
	Orderbook *OrderbookService

	TenderlyKey string
}

// NewClient creates and initializes a new Client instance based on the provided Config.
func NewClient(config Config) (*Client, error) {

	// TODO this may be replaceable with https://github.com/go-playground/validator
	err := config.validate()
	if err != nil {
		return nil, fmt.Errorf("config validation error: %v", err)
	}

	var baseUrl *url.URL
	switch config.TargetEnvironment {
	case "":
		fallthrough
	case EnvironmentProduction:
		baseUrl = baseUrlProduction
	case EnvironmentStaging:
		baseUrl = baseUrlStaging
	default:
		return nil, fmt.Errorf("unrecognized environment: %s", config.TargetEnvironment)
	}

	publicAddress := common.HexToAddress("0x0")
	if config.WalletKey != "" {
		privateKey, err := crypto.HexToECDSA(config.WalletKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert private key: %v", err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("could not cast public key to ECDSA")
		}

		publicAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	}

	var ethClient *ethclient.Client
	if config.Web3HttpProvider != "" {
		ethClient, err = ethclient.Dial(config.Web3HttpProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to create eth client: %v", err)
		}
	}

	c := &Client{
		httpClient:    &http.Client{},
		EthClient:     ethClient,
		ChainId:       config.ChainId,
		BaseURL:       baseUrl,
		ApiKey:        config.DevPortalApiKey,
		WalletKey:     config.WalletKey,
		PublicAddress: publicAddress,
		RpcUrlWithKey: config.Web3HttpProvider,
		NonceCache:    make(map[string]uint64),
		TenderlyKey:   config.TenderlyKey,
	}

	c.common.client = c

	c.Actions = (*ActionService)(&c.common)
	c.Swap = (*SwapService)(&c.common)
	c.Orderbook = (*OrderbookService)(&c.common)

	return c, nil
}

func (c *Client) NewRequest(method, urlStr string, body []byte) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
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
