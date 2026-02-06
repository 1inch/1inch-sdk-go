package fusionplus

import (
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	geth_common "github.com/ethereum/go-ethereum/common"
)

func NewExtensionPlus(params ExtensionParamsPlus) (*ExtensionPlus, error) {
	if err := fusionorder.ValidateExtensionHexParams(fusionorder.ExtensionHexParams{
		SettlementContract: params.SettlementContract,
		MakerAssetSuffix:   params.MakerAssetSuffix,
		TakerAssetSuffix:   params.TakerAssetSuffix,
		Predicate:          params.Predicate,
		CustomData:         params.CustomData,
	}); err != nil {
		return nil, err
	}

	settlementContractAddress := geth_common.HexToAddress(params.SettlementContract)
	// FusionPlus uses encoding without point count byte
	makingAndTakingAmountData := settlementContractAddress.String() + hexadecimal.Trim0x(params.AuctionDetails.EncodeWithoutPointCount())

	extensionPlus := &ExtensionPlus{
		SettlementContract:  params.SettlementContract,
		AuctionDetails:      params.AuctionDetails,
		PostInteractionData: params.PostInteractionData,
		Asset:               params.Asset,
		Permit:              params.Permit,

		MakerAssetSuffix: params.MakerAssetSuffix,
		TakerAssetSuffix: params.TakerAssetSuffix,
		MakingAmountData: makingAndTakingAmountData,
		TakingAmountData: makingAndTakingAmountData,
		Predicate:        params.Predicate,
		PreInteraction:   params.PreInteraction,
		CustomData:       params.CustomData,
	}

	postInteractionDataEncoded, err := params.PostInteractionData.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode post interaction data: %w", err)
	}
	postInteraction, err := fusionorder.NewInteraction(settlementContractAddress, postInteractionDataEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to create post interaction: %w", err)
	}
	extensionPlus.PostInteraction = postInteraction.Encode()

	if params.Permit != "" {
		permitInteraction := &fusionorder.Interaction{
			Target: geth_common.HexToAddress(params.Asset),
			Data:   params.Permit,
		}
		extensionPlus.MakerPermit = permitInteraction.Target.String() + hexadecimal.Trim0x(permitInteraction.Data)
	}

	return extensionPlus, nil
}

// Keccak256 calculates the Keccak256 hash of the extension data
func (e *ExtensionPlus) Keccak256() (*big.Int, error) {
	return fusionorder.Keccak256Hash(e)
}

func (e *ExtensionPlus) ConvertToOrderbookExtension() *orderbook.Extension {
	return &orderbook.Extension{
		MakerAssetSuffix: e.MakerAssetSuffix,
		TakerAssetSuffix: e.TakerAssetSuffix,
		MakingAmountData: e.MakingAmountData,
		TakingAmountData: e.TakingAmountData,
		Predicate:        e.Predicate,
		MakerPermit:      e.MakerPermit,
		PreInteraction:   e.PreInteraction,
		PostInteraction:  e.PostInteraction,
		//hexadecimal.Trim0x(e.CustomData), // TODO Blocking custom data for now because it is breaking the cumsum method. The extension constructor will return with an error if the user provides this field.
	}
}

func (e *ExtensionPlus) GenerateSalt() (*big.Int, error) {
	hash, err := e.Keccak256()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate extension hash: %w", err)
	}
	return fusionorder.GenerateSaltWithExtension(hash, e.isEmpty())
}

// isEmpty checks if the extension data is empty
func (e *ExtensionPlus) isEmpty() bool {
	return *e == (ExtensionPlus{})
}

func DecodeExtension(data []byte) (*ExtensionPlus, error) {
	orderbookExtension, err := orderbook.Decode(data)
	if err != nil {
		return &ExtensionPlus{}, fmt.Errorf("failed to decode extension: %w", err)
	}

	extensionPlus, err := FromLimitOrderExtension(orderbookExtension)
	if err != nil {
		return nil, fmt.Errorf("failed to convert orderbook extension to fusionplus extension: %w", err)
	}

	return &ExtensionPlus{
		SettlementContract:  extensionPlus.SettlementContract,
		AuctionDetails:      extensionPlus.AuctionDetails,
		PostInteractionData: extensionPlus.PostInteractionData,
		Asset:               extensionPlus.Asset,
		Permit:              extensionPlus.Permit,

		MakerAssetSuffix: orderbookExtension.MakerAssetSuffix,
		TakerAssetSuffix: orderbookExtension.TakerAssetSuffix,
		MakingAmountData: orderbookExtension.MakingAmountData,
		TakingAmountData: orderbookExtension.TakingAmountData,
		Predicate:        orderbookExtension.Predicate,
		MakerPermit:      orderbookExtension.MakerPermit,
		PreInteraction:   orderbookExtension.PreInteraction,
		PostInteraction:  orderbookExtension.PostInteraction,
	}, nil
}

func FromLimitOrderExtension(extension *orderbook.Extension) (*ExtensionPlus, error) {

	settlementContractAddress := extension.MakingAmountData[:42]

	if settlementContractAddress != extension.TakingAmountData[:42] {
		return nil, fmt.Errorf("malformed extension: mismatched settlement contract in making and taking amount data")
	}
	if settlementContractAddress != extension.PostInteraction[:42] {
		return nil, fmt.Errorf("malformed extension: mismatched settlement contract in making and post interaction")
	}

	auctionDetails, err := fusionorder.DecodeAuctionDetails(extension.MakingAmountData[42:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode auction details: %w", err)
	}

	postInteractionData, err := DecodeSettlementPostInteractionData(extension.PostInteraction[42:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode post interaction data: %w", err)
	}

	extensionPlus := &ExtensionPlus{
		SettlementContract:  settlementContractAddress,
		AuctionDetails:      auctionDetails,
		PostInteractionData: postInteractionData,

		MakerAssetSuffix: extension.MakerAssetSuffix,
		TakerAssetSuffix: extension.TakerAssetSuffix,
		MakingAmountData: extension.MakingAmountData,
		TakingAmountData: extension.TakingAmountData,
		Predicate:        extension.Predicate,
		MakerPermit:      extension.MakerPermit,
		PreInteraction:   extension.PreInteraction,
		PostInteraction:  extension.PostInteraction,
	}

	var permitInteraction *fusionorder.Interaction
	if extension.MakerPermit != "" && extension.MakerPermit != "0x" {
		permitInteraction, err = fusionorder.DecodeInteraction(extension.MakerPermit)
		if err != nil {
			return nil, fmt.Errorf("failed to decode permit interaction: %w", err)
		}

		extensionPlus.Asset = permitInteraction.Target.String()
		extensionPlus.Permit = permitInteraction.Data
	}

	return extensionPlus, nil
}
