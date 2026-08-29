package fusionplus

// Compile-time assertions that the v5 renames keep their old names working as
// deprecated aliases. A broken alias fails to compile here rather than in a
// downstream project.
var (
	_ GasCostConfigClass = GasCostConfig{}
	_ AuctionPointClass  = AuctionPoint{}

	_ MyMerkleTree              = MerkleTree{}
	_ ExtensionPlus             = Extension{}
	_ GetOrderByOrderHashParams = GetOrderFillsByHashParams{}

	_ QuoterControllerGetQuoteParamsFixed = QuoteParams{}
	_ GetQuoteOutputFixed                 = Quote{}
	_ GetOrderFillsByHashOutputFixed      = GetOrderFillsByHashOutput{}
)
