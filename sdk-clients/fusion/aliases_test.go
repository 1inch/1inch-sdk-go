package fusion

// Compile-time assertions that the clean, fusionplus-consistent names and the
// deprecated *Fixed/*Class aliases all resolve to the generated quoter types.
// If a rename or codegen change breaks one of these aliases, the package fails
// to compile here — turning a would-be downstream break into a local test error.
var (
	_ Preset        = PresetClass{}
	_ QuotePresets  = QuotePresetsClass{}
	_ AuctionPoint  = AuctionPointClass{}
	_ GasCostConfig = GasCostConfigClass{}

	_ PresetClass        = Preset{}
	_ QuotePresetsClass  = QuotePresets{}
	_ AuctionPointClass  = AuctionPoint{}
	_ GasCostConfigClass = GasCostConfig{}

	// Deprecated *Fixed aliases must still resolve.
	_ PresetClassFixed       = Preset{}
	_ QuotePresetsClassFixed = QuotePresets{}

	_ FusionOrderConstructor = OrderConstructor{}
)
