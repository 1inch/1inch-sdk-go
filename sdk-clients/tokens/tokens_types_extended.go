package tokens

// Deprecated: Use ProviderTokenDto instead. The type bugs this type used to
// correct (tags being objects, not strings) are now fixed at generation time
// (codegen/overrides.go); the alias is kept so existing integrations keep
// compiling.
type ProviderTokenDtoFixed = ProviderTokenDto

// Deprecated: Use TokenInfoDto instead; see ProviderTokenDtoFixed.
type TokenInfoDtoFixed = TokenInfoDto

type CustomTokensControllerGetTokenInfoParams struct {
	Address string `url:"address" json:"address"`
}
