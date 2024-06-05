// Package tokens provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package tokens

// BadRequestErrorDto defines model for BadRequestErrorDto.
type BadRequestErrorDto struct {
	Error      string  `json:"error"`
	Message    string  `json:"message"`
	StatusCode float32 `json:"statusCode"`
}

// ProviderTokenDto defines model for ProviderTokenDto.
type ProviderTokenDto struct {
	Address         string   `json:"address"`
	ChainId         float32  `json:"chainId"`
	Decimals        float32  `json:"decimals"`
	DisplayedSymbol *string  `json:"displayedSymbol,omitempty"`
	Eip2612         *bool    `json:"eip2612,omitempty"`
	IsFoT           *bool    `json:"isFoT,omitempty"`
	LogoURI         *string  `json:"logoURI,omitempty"`
	Name            string   `json:"name"`
	Providers       []string `json:"providers"`
	Symbol          string   `json:"symbol"`
	Tags            []string `json:"tags"`
}

// TagDto defines model for TagDto.
type TagDto struct {
	Provider string `json:"provider"`
	Value    string `json:"value"`
}

// TokenDto defines model for TokenDto.
type TokenDto struct {
	Address   string   `json:"address"`
	ChainId   float32  `json:"chainId"`
	Decimals  float32  `json:"decimals"`
	Eip2612   *bool    `json:"eip2612,omitempty"`
	IsFoT     *bool    `json:"isFoT,omitempty"`
	LogoURI   *string  `json:"logoURI,omitempty"`
	Name      string   `json:"name"`
	Providers []string `json:"providers"`
	Rating    float32  `json:"rating"`
	Symbol    string   `json:"symbol"`
	Tags      []TagDto `json:"tags"`
}

// TokenInfoDto defines model for TokenInfoDto.
type TokenInfoDto struct {
	Address    string                  `json:"address"`
	ChainId    float32                 `json:"chainId"`
	Decimals   float32                 `json:"decimals"`
	Extensions *map[string]interface{} `json:"extensions,omitempty"`
	LogoURI    string                  `json:"logoURI"`
	Name       string                  `json:"name"`
	Symbol     string                  `json:"symbol"`
	Tags       []string                `json:"tags"`
}

// TokenListResponseDto defines model for TokenListResponseDto.
type TokenListResponseDto struct {
	Keywords  []string          `json:"keywords"`
	LogoURI   string            `json:"logoURI"`
	Name      string            `json:"name"`
	Tags      map[string]TagDto `json:"tags"`
	TagsOrder []string          `json:"tags_order"`
	Timestamp string            `json:"timestamp"`
	Tokens    []TokenInfoDto    `json:"tokens"`
	Version   VersionDto        `json:"version"`
}

// VersionDto defines model for VersionDto.
type VersionDto struct {
	Major float32 `json:"major"`
	Minor float32 `json:"minor"`
	Patch float32 `json:"patch"`
}

// SearchControllerSearchAllChainsParams defines parameters for SearchControllerSearchAllChains.
type SearchControllerSearchAllChainsParams struct {
	// Query Text to search for in token address, token symbol, or description
	Query *string `url:"query,omitempty" json:"query,omitempty"`

	// IgnoreListed Whether to ignore listed tokens
	IgnoreListed       *bool `url:"ignore_listed,omitempty" json:"ignore_listed,omitempty"`
	OnlyPositiveRating bool  `url:"only_positive_rating" json:"only_positive_rating"`

	// Limit Maximum number of tokens to return
	Limit *float32 `url:"limit,omitempty" json:"limit,omitempty"`
}

// TokenListControllerTokensParams defines parameters for TokenListControllerTokens.
type TokenListControllerTokensParams struct {
	// Provider Provider code. Default value is 1inch
	Provider *string `url:"provider,omitempty" json:"provider,omitempty"`

	// Country Country code
	Country *string `url:"country,omitempty" json:"country,omitempty"`
}

// CustomTokensControllerGetTokensInfoParams defines parameters for CustomTokensControllerGetTokensInfo.
type CustomTokensControllerGetTokensInfoParams struct {
	Addresses []string `url:"addresses" json:"addresses"`
}

// SearchControllerSearchSingleChainParams defines parameters for SearchControllerSearchSingleChain.
type SearchControllerSearchSingleChainParams struct {
	// Query Text to search for in token address, token symbol, or description
	Query *string `url:"query,omitempty" json:"query,omitempty"`

	// IgnoreListed Whether to ignore listed tokens
	IgnoreListed       *bool `url:"ignore_listed,omitempty" json:"ignore_listed,omitempty"`
	OnlyPositiveRating bool  `url:"only_positive_rating" json:"only_positive_rating"`

	// Limit Maximum number of tokens to return
	Limit *float32 `url:"limit,omitempty" json:"limit,omitempty"`
}

// TokenListControllerTokensListParams defines parameters for TokenListControllerTokensList.
type TokenListControllerTokensListParams struct {
	// Provider Provider code. Default value is "1inch"
	Provider *string `url:"provider,omitempty" json:"provider,omitempty"`

	// Country Country code
	Country *string `url:"country,omitempty" json:"country,omitempty"`
}