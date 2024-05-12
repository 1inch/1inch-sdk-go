// Package nft provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package nft

import (
	"time"
)

// Defines values for AssetProvider.
const (
	OPENSEA AssetProvider = "OPENSEA"
	POAP    AssetProvider = "POAP"
	RARIBLE AssetProvider = "RARIBLE"
)

// Defines values for GetNftsByAddressParamsChainIds.
const (
	N1             GetNftsByAddressParamsChainIds = 1
	N10            GetNftsByAddressParamsChainIds = 10
	N100           GetNftsByAddressParamsChainIds = 100
	N1313161554e09 GetNftsByAddressParamsChainIds = 1.313161554e+09
	N137           GetNftsByAddressParamsChainIds = 137
	N250           GetNftsByAddressParamsChainIds = 250
	N324           GetNftsByAddressParamsChainIds = 324
	N42161         GetNftsByAddressParamsChainIds = 42161
	N43114         GetNftsByAddressParamsChainIds = 43114
	N45            GetNftsByAddressParamsChainIds = 45
	N56            GetNftsByAddressParamsChainIds = 56
	N8217          GetNftsByAddressParamsChainIds = 8217
	N8453          GetNftsByAddressParamsChainIds = 8453
)

// Asset defines model for Asset.
type Asset struct {
	AnimationOriginalUrl string        `json:"animation_original_url"`
	AnimationUrl         *string       `json:"animation_url,omitempty"`
	AssetContract        AssetContract `json:"asset_contract"`

	// BackgroundColor The background color to be displayed with the item
	BackgroundColor *map[string]interface{} `json:"background_color,omitempty"`

	// ChainId chain id of NFT
	ChainId      *int                    `json:"chainId,omitempty"`
	Collection   *Collection             `json:"collection,omitempty"`
	Creator      *map[string]interface{} `json:"creator,omitempty"`
	Decimals     *map[string]interface{} `json:"decimals,omitempty"`
	Description  string                  `json:"description"`
	ExternalLink string                  `json:"external_link"`

	// Id The token ID of the NFT
	Id                float32 `json:"id"`
	ImageOriginalUrl  *string `json:"image_original_url,omitempty"`
	ImagePreviewUrl   *string `json:"image_preview_url,omitempty"`
	ImageThumbnailUrl *string `json:"image_thumbnail_url,omitempty"`
	ImageUrl          *string `json:"image_url,omitempty"`
	IsNsfw            *bool   `json:"is_nsfw,omitempty"`

	// LastSale When this item was last sold (null if there was no last sale)
	LastSale    *map[string]interface{} `json:"last_sale,omitempty"`
	ListingDate *map[string]interface{} `json:"listing_date,omitempty"`

	// Name Name of the item
	Name      string   `json:"name"`
	NumSales  *float32 `json:"num_sales,omitempty"`
	Owner     *string  `json:"owner,omitempty"`
	Permalink *string  `json:"permalink,omitempty"`

	// Priority priority of NFT if "zero" it should display first
	Priority float32 `json:"priority"`

	// Provider provider of NFT
	Provider AssetProvider `json:"provider"`

	// RarityData Contains rarity data for the asset. See Rarity Data below.
	RarityData        *RarityData             `json:"rarity_data,omitempty"`
	SeaportSellOrders *map[string]interface{} `json:"seaport_sell_orders,omitempty"`
	SupportsWyvern    *bool                   `json:"supports_wyvern,omitempty"`

	// TokenId The token ID of the NFT
	TokenId       *string                 `json:"token_id,omitempty"`
	TokenMetadata *string                 `json:"token_metadata,omitempty"`
	TopBid        *map[string]interface{} `json:"top_bid,omitempty"`

	// Traits A list of traits associated with the item (see traits section)
	Traits                  []string                `json:"traits"`
	TransferFee             *map[string]interface{} `json:"transfer_fee,omitempty"`
	TransferFeePaymentToken *map[string]interface{} `json:"transfer_fee_payment_token,omitempty"`
}

// AssetProvider provider of NFT
type AssetProvider string

// AssetContract defines model for AssetContract.
type AssetContract struct {
	// Address on chain address of the contract
	Address string `json:"address"`

	// AssetContractType describes whether a contract is fungible or non-fungible
	AssetContractType string    `json:"asset_contract_type"`
	CreatedDate       time.Time `json:"created_date"`

	// Description 	description of the contract
	Description string `json:"description"`

	// ExternalLink external link to the contracts website
	ExternalLink string `json:"external_link"`

	// ImageUrl An image for the item. Note that this is the cached URL we store on our end. The original image url is image_original_url
	ImageUrl string `json:"image_url"`

	// Name name of the contract
	Name           string                 `json:"name"`
	NftVersion     string                 `json:"nft_version"`
	OpenseaVersion map[string]interface{} `json:"opensea_version"`
	Owner          string                 `json:"owner"`

	// SchemaName types of tokens supported by the contract (ex. ERC721)
	SchemaName  string `json:"schema_name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"total_supply"`
}

// AssetsResponse defines model for AssetsResponse.
type AssetsResponse struct {
	Assets Asset `json:"assets"`
}

// Collection defines model for Collection.
type Collection struct {
	// BannerImageUrl Image used in the horizontal top banner for the collection
	BannerImageUrl string `json:"banner_image_url"`

	// Description Description for the model
	Description string `json:"description"`
	DiscordUrl  string `json:"discord_url"`

	// ExternalUrl External link to the original website for the collection
	ExternalUrl string   `json:"external_url"`
	Fees        []string `json:"fees"`

	// ImageUrl An image for the collection. Note that this is the cached URL we store on our end. The original image url is image_original_url
	ImageUrl          string `json:"image_url"`
	InstagramUsername string `json:"instagram_username"`

	// Name The collection name. Typically derived from the first contract imported to the collection but can be changed by the user
	Name            string                 `json:"name"`
	Slug            map[string]interface{} `json:"slug"`
	TelegramUrl     string                 `json:"telegram_url"`
	TwitterUsername string                 `json:"twitter_username"`
	WikiUrl         map[string]interface{} `json:"wiki_url"`
}

// RarityData defines model for RarityData.
type RarityData struct {
	// CalculatedAt The time we calculated rarity data at, as a timestamp in UTC. Note: This may not equal the time a creator has uploaded or changed metadata.
	CalculatedAt time.Time `json:"calculated_at"`

	// MaxRank The maximum rank in the collection. Ranking for an asset should be considered the out of <max_rank>. Typically max_rank will be equal to collection supply if all assets have been fully revealed with metadata loaded into opensea system, and every asset has unique ranks.  Before reveal is complete, rarity_data.max_rank < collection.stats.total_supply.
	MaxRank float32 `json:"max_rank"`

	// Rank 	The rank of the asset within the collection, calculated using the rarity strategy defined by strategy_id and strategy_version.
	Rank float32 `json:"rank"`

	// RankingFeatures A dictionary of other asset features that impact rarity ranking, as returned by OpenRarity.  Currently only has "unique_attribute_count" field, which contains the number of unique attributes the asset has.
	RankingFeatures map[string]interface{} `json:"ranking_features"`

	// Score The rarity score of the asset, calculated using the rarity strategy defined by strategy_id and strategy_version.
	Score float32 `json:"score"`

	// StrategyId The rarity strategy string identifier. Current value will be "openrarity”.
	StrategyId string `json:"strategy_id"`

	// StrategyVersion The version of the strategy. For “openrarity”, this will be the python package version of the OpenRarity library used to calculate the returned score and rank.
	StrategyVersion string `json:"strategy_version"`

	// TokensScored The total tokens in the collection that have non-null traits and was used to calculate rarity data.  This will equal collection.stats.total_supply once Opensea has the revealed trait data for all assets in the collection.
	TokensScored float32 `json:"tokens_scored"`
}

// GetNftsByAddressParams defines parameters for GetNftsByAddress.
type GetNftsByAddressParams struct {
	// ChainIds List of chainIds, right now supported only ethereum & gnosis
	ChainIds []GetNftsByAddressParamsChainIds `url:"chainIds" json:"chainIds"`

	// Address web3 address of owner NFTS
	Address string `url:"address" json:"address"`

	// Limit The maximum number of api to return
	Limit *int `url:"limit,omitempty" json:"limit,omitempty"`

	// Offset The offset number of api to return
	Offset *int `url:"offset,omitempty" json:"offset,omitempty"`
}

// GetNftsByAddressParamsChainIds defines parameters for GetNftsByAddress.
type GetNftsByAddressParamsChainIds int
