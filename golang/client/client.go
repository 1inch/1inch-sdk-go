package client

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/go-querystring/query"

	"1inch-sdk-golang/helpers"
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
	TargetEnvironment          Environment
	ChainId                    int
	DevPortalApiKey            string
	Web3HttpProviderUrlWithKey string
	EtherscanApiKey            string
	WalletAddress              string
	WalletKey                  string
	LimitOrderContract         string // TODO Probably want to move this somewhere else
}

func (c *Config) validate() error {

	if c.DevPortalApiKey == "" {
		return fmt.Errorf("API key is required")
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
	// A struct that will contain a reference to this client. Used to separate each API into a unique namespace to aid in method discovery
	common service
	// Isolated namespaces for each API
	Swap        *SwapService
	TokenPrices *TokenPricesService
	Orderbook   *OrderbookService
	Fusion      *FusionService
}

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

	chainId := config.ChainId
	if chainId != 0 {
		if !helpers.IsValidChainId(chainId) {
			return nil, fmt.Errorf("invalid chain id: %d", chainId)
		}
	} else {
		chainId = 1
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
	if config.Web3HttpProviderUrlWithKey != "" {
		ethClient, err = ethclient.Dial(config.Web3HttpProviderUrlWithKey) // TODO Should the user pass this in?
		if err != nil {
			log.Fatalf("Failed to create eth client: %v", err)
		}
	}

	c := &Client{
		httpClient:    &http.Client{},
		EthClient:     ethClient,
		ChainId:       chainId,
		BaseURL:       baseUrl,
		ApiKey:        config.DevPortalApiKey,
		WalletKey:     config.WalletKey,
		PublicAddress: publicAddress,
		RpcUrlWithKey: config.Web3HttpProviderUrlWithKey,
	}

	c.common.client = c

	c.Swap = (*SwapService)(&c.common)
	c.TokenPrices = (*TokenPricesService)(&c.common)
	c.Orderbook = (*OrderbookService)(&c.common)
	c.Fusion = (*FusionService)(&c.common)

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
