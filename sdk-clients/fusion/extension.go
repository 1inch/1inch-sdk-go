package fusion

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	geth_common "github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

// TODO find a better home for these
var BASE_1E5 = big.NewInt(100000)
var BASE_1E2 = big.NewInt(100)

// Extension represents the extension data for the Fusion order
// and should be only created using the NewExtension function
type Extension struct {
	// Raw unencoded data
	SettlementContract         string
	AuctionDetails             *AuctionDetails
	PostInteractionData        *SettlementPostInteractionData
	PostInteractionDataEncoded string
	Asset                      string
	Permit                     string
	Fees                       *FeesNew
	Surplus                    *SurplusParams
	ResolvingStartTime         *big.Int

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
	SettlementContract         string
	AuctionDetails             *AuctionDetails
	PostInteractionData        *SettlementPostInteractionData
	PostInteractionDataEncoded string
	Asset                      string
	Permit                     string
	Surplus                    *SurplusParams
	ResolvingStartTime         *big.Int

	MakerAssetSuffix string
	TakerAssetSuffix string
	Predicate        string
	PreInteraction   string
	CustomData       string
}

func prefix0x(value string) string {
	if len(value) >= 2 && value[:2] == "0x" {
		return value
	}
	return "0x" + value
}

func NewExtension(params ExtensionParams) (*Extension, error) {
	if !hexadecimal.IsHexBytes(params.SettlementContract) {
		return nil, errors.New("Settlement contract must be valid hex string")
	}
	if !hexadecimal.IsHexBytes(params.MakerAssetSuffix) {
		return nil, errors.New("MakerAssetSuffix must be valid hex string")
	}
	if !hexadecimal.IsHexBytes(params.TakerAssetSuffix) {
		return nil, errors.New("TakerAssetSuffix must be valid hex string")
	}
	if !hexadecimal.IsHexBytes(params.Predicate) {
		return nil, errors.New("Predicate must be valid hex string")
	}
	if params.CustomData != "" {
		return nil, errors.New("CustomData is not currently supported")
	}

	bagdParams := &BuildAmountGetterDataParams{
		AuctionDetails:      params.AuctionDetails,
		PostInteractionData: params.PostInteractionData,
		ResolvingStartTime:  params.ResolvingStartTime,
	}

	amountData, err := BuildAmountGetterData(bagdParams, true)
	if err != nil {
		return nil, fmt.Errorf("failed to build amount getter data: %v", err)
	}

	settlementContractAddress := geth_common.HexToAddress(params.SettlementContract)
	//makingAndTakingAmountData := settlementContractAddress.String() + hexadecimal.Trim0x(params.AuctionDetails.Encode())
	makingAndTakingAmountData := strings.ToLower(settlementContractAddress.String()) + hexadecimal.Trim0x(amountData)

	fusionExtension := &Extension{
		SettlementContract:  params.SettlementContract,
		AuctionDetails:      params.AuctionDetails,
		PostInteractionData: params.PostInteractionData,
		Asset:               params.Asset,
		Permit:              params.Permit,
		Surplus:             params.Surplus,
		ResolvingStartTime:  params.ResolvingStartTime,

		MakerAssetSuffix: prefix0x(params.MakerAssetSuffix),
		TakerAssetSuffix: prefix0x(params.TakerAssetSuffix),
		MakingAmountData: prefix0x(makingAndTakingAmountData),
		TakingAmountData: prefix0x(makingAndTakingAmountData),
		Predicate:        prefix0x(params.Predicate),
		MakerPermit:      prefix0x(params.Permit),
		PreInteraction:   prefix0x(params.PreInteraction),
		CustomData:       prefix0x(params.CustomData),
	}

	postInteractionDataEncoded, err := CreateEncodedPostInteractionData(fusionExtension)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoded post interaction data: %v", err)
	}

	fusionExtension.PostInteraction = NewInteraction(settlementContractAddress, postInteractionDataEncoded).Encode()

	if params.Permit != "" {
		permitInteraction := &Interaction{
			Target: geth_common.HexToAddress(params.Asset),
			Data:   params.Permit,
		}
		fusionExtension.MakerPermit = permitInteraction.Target.String() + hexadecimal.Trim0x(permitInteraction.Data)
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
		//hexadecimal.Trim0x(e.CustomData), // TODO Blocking custom data for now because it is breaking the cumsum method. The extension constructor will return with an error if the user provides this field.
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

//func DecodeExtension(data []byte) (*Extension, error) {
//	orderbookExtension, err := orderbook.Decode(data)
//	if err != nil {
//		return &Extension{}, fmt.Errorf("error decoding extension: %v", err)
//	}
//
//	fusionExtension, err := FromLimitOrderExtension(orderbookExtension)
//	if err != nil {
//		return nil, fmt.Errorf("failed to convert orderbook extension to fusion extension: %v", err)
//	}
//
//	return &Extension{
//		SettlementContract:  fusionExtension.SettlementContract,
//		AuctionDetails:      fusionExtension.AuctionDetails,
//		PostInteractionData: fusionExtension.PostInteractionData,
//		Asset:               fusionExtension.Asset,
//		Permit:              fusionExtension.Permit,
//
//		MakerAssetSuffix: orderbookExtension.MakerAssetSuffix,
//		TakerAssetSuffix: orderbookExtension.TakerAssetSuffix,
//		MakingAmountData: orderbookExtension.MakingAmountData,
//		TakingAmountData: orderbookExtension.TakingAmountData,
//		Predicate:        orderbookExtension.Predicate,
//		MakerPermit:      orderbookExtension.MakerPermit,
//		PreInteraction:   orderbookExtension.PreInteraction,
//		PostInteraction:  orderbookExtension.PostInteraction,
//	}, nil
//}

//func FromLimitOrderExtension(extension *orderbook.Extension) (*Extension, error) {
//
//	settlementContractAddress := extension.MakingAmountData[:42]
//
//	if settlementContractAddress != extension.TakingAmountData[:42] {
//		return nil, fmt.Errorf("malfomed extension: settlement contract address should be the same in making and taking amount data")
//	}
//	if settlementContractAddress != extension.PostInteraction[:42] {
//		return nil, fmt.Errorf("malfomed extension: settlement contract address should be the same in making and post interaction")
//	}
//
//	auctionDetails, err := DecodeAuctionDetails(extension.MakingAmountData[42:])
//	if err != nil {
//		return nil, fmt.Errorf("failed to decode auction details: %v", err)
//	}
//
//	postInteractionData, err := Decode(extension.PostInteraction[42:])
//	if err != nil {
//		return nil, fmt.Errorf("failed to decode post interaction data: %v", err)
//	}
//
//	fusionExtension := &Extension{
//		SettlementContract:  settlementContractAddress,
//		AuctionDetails:      auctionDetails,
//		PostInteractionData: &postInteractionData,
//
//		MakerAssetSuffix: extension.MakerAssetSuffix,
//		TakerAssetSuffix: extension.TakerAssetSuffix,
//		MakingAmountData: extension.MakingAmountData,
//		TakingAmountData: extension.TakingAmountData,
//		Predicate:        extension.Predicate,
//		MakerPermit:      extension.MakerPermit,
//		PreInteraction:   extension.PreInteraction,
//		PostInteraction:  extension.PostInteraction,
//	}
//
//	var permitInteraction *Interaction
//	if extension.MakerPermit != "" && extension.MakerPermit != "0x" {
//		permitInteraction, err = DecodeInteraction(extension.MakerPermit)
//		if err != nil {
//			return nil, fmt.Errorf("failed to decode permit interaction: %v", err)
//		}
//
//		fusionExtension.Asset = permitInteraction.Target.String()
//		fusionExtension.Permit = permitInteraction.Data
//	}
//
//	return fusionExtension, nil
//}

type BuildAmountGetterDataParams struct {
	AuctionDetails      *AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	ResolvingStartTime  *big.Int
}

func BuildAmountGetterData(params *BuildAmountGetterDataParams, forAmountGetters bool) (string, error) {
	bytes := bytesbuilder.New()

	if forAmountGetters {
		err := bytes.AddBytes(params.AuctionDetails.Encode())
		if err != nil {
			return "", fmt.Errorf("failed to add auction details: %v", err)
		}
	}

	fee := big.NewInt(0)
	if params.PostInteractionData.AuctionFees != nil && params.PostInteractionData.AuctionFees.Integrator.Fee != nil && !params.PostInteractionData.AuctionFees.Integrator.Fee.IsZero() {
		fee = params.PostInteractionData.AuctionFees.Integrator.Fee.ToFraction(BASE_1E5)
	}
	bytes.AddUint16(fee)

	share := big.NewInt(0)
	if params.PostInteractionData.AuctionFees != nil && params.PostInteractionData.AuctionFees.Integrator.Share != nil && !params.PostInteractionData.AuctionFees.Integrator.Share.IsZero() {
		share = params.PostInteractionData.AuctionFees.Integrator.Share.ToFraction(BASE_1E2)
	}
	bytes.AddUint8(uint8(share.Uint64()))

	resolverFee := big.NewInt(0)
	if params.PostInteractionData.AuctionFees != nil && params.PostInteractionData.AuctionFees.Resolver.Fee != nil && !params.PostInteractionData.AuctionFees.Resolver.Fee.IsZero() {
		resolverFee = params.PostInteractionData.AuctionFees.Resolver.Fee.ToFraction(BASE_1E5)
	}
	bytes.AddUint16(resolverFee)

	whitelistDiscount := BpsZero
	if params.PostInteractionData.AuctionFees != nil && params.PostInteractionData.AuctionFees.Resolver.Fee != nil && !params.PostInteractionData.AuctionFees.Resolver.Fee.IsZero() {
		whitelistDiscount = params.PostInteractionData.AuctionFees.Resolver.WhitelistDiscount
	}
	discountValue := whitelistDiscount.ToFraction(BASE_1E2)
	discountNumerator := new(big.Int).Sub(new(big.Int).Set(BASE_1E2), discountValue)

	bytes.AddUint8(uint8(discountNumerator.Uint64()))

	// TODO find whitelist, add it, then write unit tests to match javascript

	if forAmountGetters {
		// Whitelist address halves only, no delays
		numWhitelist := len(params.PostInteractionData.Whitelist)
		bytes.AddUint8(uint8(numWhitelist))

		for _, entry := range params.PostInteractionData.Whitelist {
			err := bytes.AddBytes(entry.AddressHalf)
			if err != nil {
				return "", fmt.Errorf("failed to add whitelist address half: %v", err)
			}
		}
	} else {
		// Add resolvingStartTime as uint32
		bytes.AddUint32(params.ResolvingStartTime)

		// Add whitelist length as uint8
		bytes.AddUint8(uint8(len(params.PostInteractionData.Whitelist)))

		// Add each whitelist entry: addressHalf + delay
		for _, entry := range params.PostInteractionData.Whitelist {
			if err := bytes.AddBytes(fmt.Sprintf("0x%s", entry.AddressHalf)); err != nil {
				return "", fmt.Errorf("failed to add addressHalf: %w", err)
			}
			bytes.AddUint16(entry.Delay)
		}
	}

	return fmt.Sprintf("0x%s", bytes.AsHex()), nil
}
