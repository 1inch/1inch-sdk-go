package fusionorder

// OrderInfo contains the core order information shared between Fusion and FusionPlus
type OrderInfo struct {
	// Maker is the address of the account creating the order
	Maker string `json:"maker"`

	// MakerAsset is the identifier of the asset being offered by the maker
	MakerAsset string `json:"makerAsset"`

	// MakingAmount is the amount of the makerAsset being offered
	MakingAmount string `json:"makingAmount"`

	// Receiver is the address that will receive the taker asset
	Receiver string `json:"receiver"`

	// TakerAsset is the identifier of the asset being requested
	TakerAsset string `json:"takerAsset"`

	// TakingAmount is the amount of the takerAsset being requested
	TakingAmount string `json:"takingAmount"`
}

// NewOrderInfo creates a new OrderInfo
func NewOrderInfo(maker, makerAsset, makingAmount, receiver, takerAsset, takingAmount string) *OrderInfo {
	return &OrderInfo{
		Maker:        maker,
		MakerAsset:   makerAsset,
		MakingAmount: makingAmount,
		Receiver:     receiver,
		TakerAsset:   takerAsset,
		TakingAmount: takingAmount,
	}
}
