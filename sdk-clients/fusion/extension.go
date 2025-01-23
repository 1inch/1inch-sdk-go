package fusion

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/internal/hexidecimal"
	geth_common "github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

// Extension represents the extension data for the Fusion order
// and should be only created using the NewExtension function
type Extension struct {
	// Raw unencoded data
	SettlementContract  string
	AuctionDetails      *AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Asset               string
	Permit              string

	// Data formatted for Limit Order Extension
	MakerAssetSuffix string
	TakerAssetSuffix string
	MakingAmountData string
	TakingAmountData string
	Predicate        string
	MakerPermit      string
	PreInteraction   string
	PostInteraction  string
	CustomData       string
}

type ExtensionParams struct {
	SettlementContract  string
	AuctionDetails      *AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Asset               string
	Permit              string

	MakerAssetSuffix string
	TakerAssetSuffix string
	Predicate        string
	PreInteraction   string
	CustomData       string
}

func NewExtension(params ExtensionParams) (*Extension, error) {
	if !hexidecimal.IsHexBytes(params.SettlementContract) {
		return nil, errors.New("Settlement contract must be valid hex string")
	}
	if !hexidecimal.IsHexBytes(params.MakerAssetSuffix) {
		return nil, errors.New("MakerAssetSuffix must be valid hex string")
	}
	if !hexidecimal.IsHexBytes(params.TakerAssetSuffix) {
		return nil, errors.New("TakerAssetSuffix must be valid hex string")
	}
	if !hexidecimal.IsHexBytes(params.Predicate) {
		return nil, errors.New("Predicate must be valid hex string")
	}
	if params.CustomData != "" {
		return nil, errors.New("CustomData is not currently supported")
	}

	settlementContractAddress := geth_common.HexToAddress(params.SettlementContract)
	makingAndTakingAmountData := settlementContractAddress.String() + hexidecimal.Trim0x(params.AuctionDetails.Encode())

	fusionExtension := &Extension{
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

	postInteractoinDataEncoded, err := params.PostInteractionData.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode post interaction data: %v", err)
	}
	fusionExtension.PostInteraction = NewInteraction(settlementContractAddress, postInteractoinDataEncoded).Encode()

	if params.Permit != "" {
		permitInteraction := &Interaction{
			Target: geth_common.HexToAddress(params.Asset),
			Data:   params.Permit,
		}
		fusionExtension.MakerPermit = permitInteraction.Target.String() + hexidecimal.Trim0x(permitInteraction.Data)
	}

	return fusionExtension, nil
}

// Keccak256 calculates the Keccak256 hash of the extension data
func (e *Extension) Keccak256() *big.Int {
	jsonData, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	hash := sha3.New256()
	hash.Write(jsonData)
	return new(big.Int).SetBytes(hash.Sum(nil))
}

func (e *Extension) ConvertToOrderbookExtension() *orderbook.Extension {
	return &orderbook.Extension{
		MakerAssetSuffix: e.MakerAssetSuffix,
		TakerAssetSuffix: e.TakerAssetSuffix,
		MakingAmountData: e.MakingAmountData,
		TakingAmountData: e.TakingAmountData,
		Predicate:        e.Predicate,
		MakerPermit:      e.MakerPermit,
		PreInteraction:   e.PreInteraction,
		PostInteraction:  e.PostInteraction,
		//hexidecimal.Trim0x(e.CustomData), // TODO Blocking custom data for now because it is breaking the cumsum method. The extension constructor will return with an error if the user provides this field.
	}
}

func (e *Extension) GenerateSalt() (*big.Int, error) {

	// Define the maximum value (2^96 - 1)
	maxValue := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 96), big.NewInt(1))

	// Generate a random big.Int within the range [0, 2^96 - 1]
	baseSalt, err := random_number_generation.BigIntMaxFunc(maxValue)
	if err != nil {
		return nil, err
	}

	if e.isEmpty() {
		return baseSalt, nil
	}

	uint160Max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))

	extensionHash := e.Keccak256()
	salt := new(big.Int).Lsh(baseSalt, 160)
	salt.Or(salt, new(big.Int).And(extensionHash, uint160Max))

	return salt, nil
}

// isEmpty checks if the extension data is empty
func (e *Extension) isEmpty() bool {
	return *e == (Extension{})
}

func DecodeExtension(data []byte) (*Extension, error) {
	orderbookExtension, err := orderbook.Decode(data)
	if err != nil {
		return &Extension{}, fmt.Errorf("error decoding extension: %v", err)
	}

	fusionExtension, err := FromLimitOrderExtension(orderbookExtension)
	if err != nil {
		return nil, fmt.Errorf("failed to convert orderbook extension to fusion extension: %v", err)
	}

	return &Extension{
		SettlementContract:  fusionExtension.SettlementContract,
		AuctionDetails:      fusionExtension.AuctionDetails,
		PostInteractionData: fusionExtension.PostInteractionData,
		Asset:               fusionExtension.Asset,
		Permit:              fusionExtension.Permit,

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

func FromLimitOrderExtension(extension *orderbook.Extension) (*Extension, error) {

	settlementContractAddress := extension.MakingAmountData[:42]

	if settlementContractAddress != extension.TakingAmountData[:42] {
		return nil, fmt.Errorf("malfomed extension: settlement contract address should be the same in making and taking amount data")
	}
	if settlementContractAddress != extension.PostInteraction[:42] {
		return nil, fmt.Errorf("malfomed extension: settlement contract address should be the same in making and post interaction")
	}

	auctionDetails, err := DecodeAuctionDetails(extension.MakingAmountData[42:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode auction details: %v", err)
	}

	postInteractionData, err := Decode(extension.PostInteraction[42:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode post interaction data: %v", err)
	}

	fusionExtension := &Extension{
		SettlementContract:  fmt.Sprintf("%s", settlementContractAddress),
		AuctionDetails:      auctionDetails,
		PostInteractionData: &postInteractionData,

		MakerAssetSuffix: extension.MakerAssetSuffix,
		TakerAssetSuffix: extension.TakerAssetSuffix,
		MakingAmountData: extension.MakingAmountData,
		TakingAmountData: extension.TakingAmountData,
		Predicate:        extension.Predicate,
		MakerPermit:      extension.MakerPermit,
		PreInteraction:   extension.PreInteraction,
		PostInteraction:  extension.PostInteraction,
	}

	var permitInteraction *Interaction
	if extension.MakerPermit != "" && extension.MakerPermit != "0x" {
		permitInteraction, err = DecodeInteraction(extension.MakerPermit)
		if err != nil {
			return nil, fmt.Errorf("failed to decode permit interaction: %v", err)
		}

		fusionExtension.Asset = permitInteraction.Target.String()
		fusionExtension.Permit = permitInteraction.Data
	}

	return fusionExtension, nil
}
