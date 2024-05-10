package orderbook

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	thresholdMask             = NewBitMask(big.NewInt(0), big.NewInt(185))
	argsInteractionLengthMask = NewBitMask(big.NewInt(220), big.NewInt(224))
	argsExtensionLengthMask   = NewBitMask(big.NewInt(224), big.NewInt(248))
)

const (
	MakerAmountFlag     = 255
	UnwrapWETHFlag      = 254
	SkipOrderPermitFlag = 253
	UsePermit2Flag      = 252
	ArgsHasReceiverFlag = 251
)

type TakerTraitsEncoded struct {
	TraitFlags *big.Int
	Args       []byte
}

type TakerTraits struct {
	Receiver  *common.Address
	Extension string // Assuming extension related functions are defined elsewhere

	MakerAmount     bool
	UnwrapWETH      bool
	SkipOrderPermit bool
	UsePermit2      bool
	ArgsHasReceiver bool
}

func NewTakerTraits(params TakerTraitsParams) *TakerTraits {
	return &TakerTraits{
		Receiver:        params.Receiver,
		Extension:       params.Extension,
		MakerAmount:     params.MakerAmount,
		UnwrapWETH:      params.UnwrapWETH,
		SkipOrderPermit: params.SkipOrderPermit,
		UsePermit2:      params.UsePermit2,
		ArgsHasReceiver: params.ArgsHasReceiver,
	}
}

func (t *TakerTraits) Encode() *TakerTraitsEncoded {
	encodedCalldata := new(big.Int)
	tmp := new(big.Int)

	if t.ArgsHasReceiver && t.Receiver != nil {
		encodedCalldata.Or(encodedCalldata, tmp.Lsh(big.NewInt(1), ArgsHasReceiverFlag))
	}

	var extensionBytesLen *big.Int
	if t.Extension != "0x" {
		extensionBytesLen = big.NewInt(int64(len(t.Extension))/2 - 1)
		argsExtensionLengthMask.SetBits(encodedCalldata, extensionBytesLen)
	} else {
		extensionBytesLen = big.NewInt(0)
	}

	traits := fmt.Sprintf("%032x", encodedCalldata)
	traitsBig, err := hexutil.DecodeBig("0x" + traits)
	if err != nil {
		panic(err)
	}
	return &TakerTraitsEncoded{
		TraitFlags: traitsBig,
		Args:       common.FromHex(t.Extension),
	}
}
