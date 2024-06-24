package fusion

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

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

type ExtensionBuilder struct {
	makerAssetSuffix string
	takerAssetSuffix string
	makingAmountData string
	takingAmountData string
	predicate        string
	makerPermit      string
	preInteraction   string
	postInteraction  string
	customData       string
}

func trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func (b *ExtensionBuilder) WithMakerAssetSuffix(suffix string) *ExtensionBuilder {
	if !isHexBytes(suffix) {
		panic("MakerAssetSuffix must be valid hex string")
	}
	b.makerAssetSuffix = suffix
	return b
}

func (b *ExtensionBuilder) WithTakerAssetSuffix(suffix string) *ExtensionBuilder {
	if !isHexBytes(suffix) {
		panic("TakerAssetSuffix must be valid hex string")
	}
	b.takerAssetSuffix = suffix
	return b
}

func (b *ExtensionBuilder) WithMakingAmountData(address common.Address, data string) *ExtensionBuilder {
	if !isHexBytes(data) {
		panic("MakingAmountData must be valid hex string")
	}
	b.makingAmountData = address.String() + trim0x(data)
	return b
}

func (b *ExtensionBuilder) WithTakingAmountData(address common.Address, data string) *ExtensionBuilder {
	if !isHexBytes(data) {
		panic("TakingAmountData must be valid hex string")
	}
	b.takingAmountData = address.String() + trim0x(data)
	return b
}

func (b *ExtensionBuilder) WithPredicate(predicate string) *ExtensionBuilder {
	if !isHexBytes(predicate) {
		panic("Predicate must be valid hex string")
	}
	b.predicate = predicate
	return b
}

func (b *ExtensionBuilder) WithMakerPermit(tokenFrom common.Address, permitData string) *ExtensionBuilder {
	if !isHexBytes(permitData) {
		panic("Permit data must be valid hex string")
	}
	b.makerPermit = tokenFrom.String() + trim0x(permitData)
	return b
}

func (b *ExtensionBuilder) WithPreInteraction(interaction Interaction) *ExtensionBuilder {
	b.preInteraction = interaction.Encode()
	return b
}

func (b *ExtensionBuilder) WithPostInteraction(interaction Interaction) *ExtensionBuilder {
	b.postInteraction = interaction.Encode()
	return b
}

func (b *ExtensionBuilder) WithCustomData(data string) *ExtensionBuilder {
	if !isHexBytes(data) {
		panic("Custom data must be valid hex string")
	}
	b.customData = data
	return b
}

func (b *ExtensionBuilder) Build() *Extension {
	return &Extension{
		MakerAssetSuffix: b.makerAssetSuffix,
		TakerAssetSuffix: b.takerAssetSuffix,
		MakingAmountData: b.makingAmountData,
		TakingAmountData: b.takingAmountData,
		Predicate:        b.predicate,
		MakerPermit:      b.makerPermit,
		PreInteraction:   b.preInteraction,
		PostInteraction:  b.postInteraction,
		CustomData:       b.customData,
	}
}

// isEmpty checks if the extension data is empty
func (e *Extension) isEmpty() bool {
	return *e == (Extension{})
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
			//strings.TrimPrefix(e.CustomData, "0x"), // TODO blocking custom data for now because it is breaking the cumsum method
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
