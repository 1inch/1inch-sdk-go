package orderbook

import (
	"fmt"
	"math/big"
)

var (
// TODO currently unused masks carried over from the Typescript Limit Order SDK
// allowedSenderMask = NewBitMask(big.NewInt(0), big.NewInt(80))
// expirationMask    = NewBitMask(big.NewInt(80), big.NewInt(120))
// nonceOrEpochMask  = NewBitMask(big.NewInt(120), big.NewInt(160))
// seriesMask        = NewBitMask(big.NewInt(160), big.NewInt(200))
)

const (
	noPartialFillsFlag      = 255
	allowMultipleFillsFlag  = 254
	needPreinteractionFlag  = 252
	needPostinteractionFlag = 251
	needEpochCheckFlag      = 250
	hasExtensionFlag        = 249
	usePermit2Flag          = 248
	unwrapWethFlag          = 247
)

type MakerTraits struct {
	AllowedSender string
	Expiry        int64
	Nonce         int64
	Series        int64

	NoPartialFills      bool
	NeedPostinteraction bool
	NeedPreinteraction  bool
	NeedEpochCheck      bool
	HasExtension        bool
	ShouldUsePermit2    bool
	ShouldUnwrapWeth    bool
}

func NewMakerTraits(params MakerTraitsParams) *MakerTraits {
	return &MakerTraits{
		AllowedSender: params.AllowedSender,
		Expiry:        params.Expiry,
		Nonce:         params.Nonce,
		Series:        params.Series,

		NeedPostinteraction: params.HasPostInteraction,
		NeedPreinteraction:  params.HasPreInteraction,
		NeedEpochCheck:      params.ShouldCheckEpoch,
		HasExtension:        params.HasExtension,
		ShouldUsePermit2:    params.UsePermit2,
		ShouldUnwrapWeth:    params.UnwrapWeth,
	}
}

func (m *MakerTraits) Encode() string {
	encodedCalldata := new(big.Int)

	tmp := new(big.Int)

	if m.NoPartialFills {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), noPartialFillsFlag))
	}
	// Limit Orders require this flag to always be present
	encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), allowMultipleFillsFlag))
	if m.NeedPreinteraction {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), needPreinteractionFlag))
	}
	if m.NeedPostinteraction {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), needPostinteractionFlag))
	}
	if m.NeedEpochCheck {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), needEpochCheckFlag))
	}
	if m.HasExtension {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), hasExtensionFlag))
	}
	if m.ShouldUsePermit2 {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), usePermit2Flag))
	}
	if m.ShouldUnwrapWeth {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), unwrapWethFlag))
	}

	// Convert AllowedSender from hex string to big.Int
	allowedSenderInt := new(big.Int)
	allowedSenderInt.SetString(m.AllowedSender[len(m.AllowedSender)-20:], 16) // We only care about the last 20 characters of the ethereum address

	// TODO These values originally used masks to write. Needs more testing to verify the simpler approach works. See https://github.com/1inch/limit-order-sdk/blob/0724227f6dab1649c4a4abcb1df30c2b43126eab/src/limit-order/maker-traits.ts#L74-L84 for how this looks in the Typescript Limit Order SDK
	encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(m.Series), 160))
	encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(m.Nonce), 120))
	encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(m.Expiry), 80))
	encodedCalldata.Or(encodedCalldata, tmp.And(allowedSenderInt, new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 80), big.NewInt(1))))

	// Pad the predicate to 32 bytes with 0's on the left and convert to hex string
	paddedPredicate := fmt.Sprintf("%032x", encodedCalldata)
	return "0x" + paddedPredicate
}
