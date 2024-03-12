package models

import "fmt"

type Config struct {
	DevPortalApiKey   string
	Web3HttpProviders []Web3Provider
}

type Web3Provider struct {
	ChainId int
	Url     string
}

func (c *Config) Validate() error {

	if c.DevPortalApiKey == "" {
		return fmt.Errorf("API key is required")
	}
	if len(c.Web3HttpProviders) == 0 {
		return fmt.Errorf("at least one web3 provider URL is required")
	}
	for _, provider := range c.Web3HttpProviders {
		if provider.ChainId == 0 {
			return fmt.Errorf("all web3 providers must have a chain ID set")
		}
		if provider.Url == "" {
			return fmt.Errorf("all web3 providers must have a URL set")
		}
	}

	return nil
}
