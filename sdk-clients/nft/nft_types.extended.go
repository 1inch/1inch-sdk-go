package nft

type SupportedChainsResponse []GetNftsByAddressParamsChainIds

type Trait struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type AssetContractExtended struct {
	Address           string `json:"address"`
	AssetContractType string `json:"asset_contract_type"`
	CreatedDate       string `json:"created_date"`
	Description       string `json:"description"`
	Name              string `json:"name"`
	ImageURL          string `json:"image_url"`
	ExternalLink      string `json:"external_link"`
	TotalSupply       string `json:"total_supply"`
	Owner             string `json:"owner"`
	OpenseaVersion    string `json:"opensea_version"`
	NFTVersion        string `json:"nft_version"`
	SchemaName        string `json:"schema_name"`
	Symbol            string `json:"symbol"`
}

type AssetExtended struct {
	ID                   int64                 `json:"id"`
	Provider             string                `json:"provider"`
	TokenID              string                `json:"token_id"`
	AnimationOriginalURL string                `json:"animation_original_url"`
	AnimationURL         string                `json:"animation_url"`
	Description          string                `json:"description"`
	ExternalLink         string                `json:"external_link"`
	Permalink            string                `json:"permalink"`
	Name                 string                `json:"name"`
	ChainID              int                   `json:"chainId"`
	Traits               []Trait               `json:"traits"`
	Priority             int                   `json:"priority"`
	AssetContract        AssetContractExtended `json:"asset_contract"`
}

type GetNFTsByAddressResponse struct {
	Assets []AssetExtended `json:"assets"`
}
