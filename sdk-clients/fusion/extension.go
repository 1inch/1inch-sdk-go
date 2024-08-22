package fusion

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

// Extension represents the extension data for the Fusion order
// and should be only created using the NewExtension function
type Extension struct {
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

func NewExtension(params ExtensionParams) (*Extension, error) {
	if !isHexBytes(params.MakerAssetSuffix) {
		return nil, errors.New("MakerAssetSuffix must be valid hex string")
	}
	if !isHexBytes(params.TakerAssetSuffix) {
		return nil, errors.New("TakerAssetSuffix must be valid hex string")
	}
	if !isHexBytes(params.MakingAmountData) {
		return nil, errors.New("MakingAmountData must be valid hex string")
	}
	if !isHexBytes(params.TakingAmountData) {
		return nil, errors.New("TakingAmountData must be valid hex string")
	}
	if !isHexBytes(params.Predicate) {
		return nil, errors.New("Predicate must be valid hex string")
	}
	if !isHexBytes(params.MakerPermit) {
		return nil, errors.New("MakerPermit must be valid hex string")
	}
	if params.CustomData != "" {
		return nil, errors.New("CustomData is not currently supported")
	}
	if !isHexBytes(params.CustomData) {
		return nil, errors.New("CustomData must be valid hex string")
	}

	return &Extension{
		MakerAssetSuffix: params.MakerAssetSuffix,
		TakerAssetSuffix: params.TakerAssetSuffix,
		MakingAmountData: params.MakingAmountData,
		TakingAmountData: params.TakingAmountData,
		Predicate:        params.Predicate,
		MakerPermit:      params.MakerPermit,
		PreInteraction:   params.PreInteraction,
		PostInteraction:  params.PostInteraction,
		CustomData:       params.CustomData,
	}, nil
}

// keccak256 calculates the Keccak256 hash of the extension data
func (e *Extension) keccak256() *big.Int {
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
		InteractionsArray: []string{
			strings.TrimPrefix(e.MakerAssetSuffix, "0x"),
			strings.TrimPrefix(e.TakerAssetSuffix, "0x"),
			strings.TrimPrefix(e.MakingAmountData, "0x"),
			strings.TrimPrefix(e.TakingAmountData, "0x"),
			strings.TrimPrefix(e.Predicate, "0x"),
			strings.TrimPrefix(e.MakerPermit, "0x"),
			e.PreInteraction,
			e.PostInteraction,
			//strings.TrimPrefix(e.CustomData, "0x"), // TODO Blocking custom data for now because it is breaking the cumsum method. The extension constructor will return with an error if the user provides this field.
		},
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

	extensionHash := e.keccak256()
	salt := new(big.Int).Lsh(baseSalt, 160)
	salt.Or(salt, new(big.Int).And(extensionHash, uint160Max))

	return salt, nil
}

// isEmpty checks if the extension data is empty
func (e *Extension) isEmpty() bool {
	return *e == (Extension{})
}

func trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}
