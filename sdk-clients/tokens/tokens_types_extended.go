package tokens

// ProviderTokenDtoFixed uses the Tag struct for Tags
type ProviderTokenDtoFixed struct {
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
	Tags            []TagDto `json:"tags"`
}

// TokenInfoDtoFixed uses the Tag struct for Tags
type TokenInfoDtoFixed struct {
	Address    string                  `json:"address"`
	ChainId    float32                 `json:"chainId"`
	Decimals   float32                 `json:"decimals"`
	Extensions *map[string]any `json:"extensions,omitempty"`
	LogoURI    string                  `json:"logoURI"`
	Name       string                  `json:"name"`
	Symbol     string                  `json:"symbol"`
	Tags       []TagDto                `json:"tags"`
}

type CustomTokensControllerGetTokenInfoParams struct {
	Address string `url:"address" json:"address"`
}
