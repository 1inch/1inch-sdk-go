package fusion

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

// TODO this class design is pointless without reasonable defaults baked in. Just accept a struct and validate everything in the constructor

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

func isHexString(s string) bool {
	// Assume this function checks if a string is a valid hex string
	return true
}

func trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func (b *ExtensionBuilder) WithMakerAssetSuffix(suffix string) *ExtensionBuilder {
	if !isHexString(suffix) {
		panic("MakerAssetSuffix must be valid hex string")
	}
	b.makerAssetSuffix = suffix
	return b
}

func (b *ExtensionBuilder) WithTakerAssetSuffix(suffix string) *ExtensionBuilder {
	if !isHexString(suffix) {
		panic("TakerAssetSuffix must be valid hex string")
	}
	b.takerAssetSuffix = suffix
	return b
}

func (b *ExtensionBuilder) WithMakingAmountData(address common.Address, data string) *ExtensionBuilder {
	if !isHexString(data) {
		panic("MakingAmountData must be valid hex string")
	}
	b.makingAmountData = address.String() + trim0x(data)
	return b
}

func (b *ExtensionBuilder) WithTakingAmountData(address common.Address, data string) *ExtensionBuilder {
	if !isHexString(data) {
		panic("TakingAmountData must be valid hex string")
	}
	b.takingAmountData = address.String() + trim0x(data)
	return b
}

func (b *ExtensionBuilder) WithPredicate(predicate string) *ExtensionBuilder {
	if !isHexString(predicate) {
		panic("Predicate must be valid hex string")
	}
	b.predicate = predicate
	return b
}

func (b *ExtensionBuilder) WithMakerPermit(tokenFrom common.Address, permitData string) *ExtensionBuilder {
	if !isHexString(permitData) {
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
	if !isHexString(data) {
		panic("Custom data must be valid hex string")
	}
	b.customData = data
	return b
}

func (b *ExtensionBuilder) Build() Extension {
	return Extension{
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

func (e Extension) ConvertToOrderbookExtension() *orderbook.Extension {
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

	//var builder strings.Builder
	//for _, interaction := range extension.InteractionsArray {
	//	interaction = strings.TrimPrefix(interaction, "0x")
	//	builder.WriteString(interaction)
	//}
	//interactionsConcatenated := builder.String()
	//
	//offsetsBytes := getOffsets(extension).Bytes()
	//paddedOffsetHex := fmt.Sprintf("%064x", offsetsBytes)
	//return "0x" + paddedOffsetHex + interactionsConcatenated
}

//func getOffsets(oe orderbook.Extension) *big.Int {
//	var lengthMap []int
//	for _, interaction := range oe.InteractionsArray {
//		lengthMap = append(lengthMap, len(strings.TrimPrefix(interaction, "0x"))/2)
//	}
//
//	cumulativeSum := 0
//	bytesAccumulator := big.NewInt(0)
//	var index uint64
//
//	for _, length := range lengthMap {
//		cumulativeSum += length
//		shiftVal := big.NewInt(int64(cumulativeSum))
//		shiftVal.Lsh(shiftVal, uint(32*index))           // Shift left
//		bytesAccumulator.Add(bytesAccumulator, shiftVal) // Add to accumulator
//		index++
//	}
//
//	return bytesAccumulator
//}

var randBigIntNewFunc func(*big.Int) *big.Int = randBigIntNew

// randBigInt generates a random big.Int within the specified range
func randBigIntNew(max *big.Int) *big.Int {
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return n
}

// BuildSalt generates a salt based on the extension and a base salt
func (e *Extension) BuildSalt() *big.Int {

	// Define the maximum value (2^96 - 1)
	maxValue := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 96), big.NewInt(1))

	// Generate a random big.Int within the range [0, 2^96 - 1]
	baseSalt := randBigIntNewFunc(maxValue)

	if e.isEmpty() {
		return baseSalt
	}

	fmt.Printf("Base salt: %v\n", baseSalt)

	UINT_160_MAX := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))

	extensionHash := e.keccak256()
	salt := new(big.Int).Lsh(baseSalt, 160)
	salt.Or(salt, new(big.Int).And(extensionHash, UINT_160_MAX))

	return salt
}

func TestExtensionBuilder(t *testing.T) {
	address := common.HexToAddress("0x0000000000000000000000000000000000000000")
	interaction := Interaction{}

	builder := &ExtensionBuilder{}
	extension := builder.
		WithMakerAssetSuffix("0x1234").
		WithTakerAssetSuffix("0x5678").
		WithMakingAmountData(address, "0x9abc").
		WithTakingAmountData(address, "0xdef0").
		WithPredicate("0x1111").
		WithMakerPermit(address, "0x2222").
		WithPreInteraction(interaction).
		WithPostInteraction(interaction).
		WithCustomData("0x3333").
		Build()

	expected := Extension{
		MakerAssetSuffix: "0x1234",
		TakerAssetSuffix: "0x5678",
		MakingAmountData: "0x00000000000000000000000000000000000000009abc",
		TakingAmountData: "0x0000000000000000000000000000000000000000def0",
		Predicate:        "0x1111",
		MakerPermit:      "0x00000000000000000000000000000000000000002222",
		PreInteraction:   "encoded_interaction",
		PostInteraction:  "encoded_interaction",
		CustomData:       "0x3333",
	}

	assert.Equal(t, expected, extension)
}
