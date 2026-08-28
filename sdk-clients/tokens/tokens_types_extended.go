package tokens

// ProviderTokenDtoFixed and TokenInfoDtoFixed are kept as aliases for
// backward compatibility. The type bugs they used to correct (tags being
// objects, not strings) are now fixed at generation time in
// codegen/overrides.go, so the generated types are correct.
type ProviderTokenDtoFixed = ProviderTokenDto
type TokenInfoDtoFixed = TokenInfoDto

type CustomTokensControllerGetTokenInfoParams struct {
	Address string `url:"address" json:"address"`
}
