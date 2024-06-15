package orderbook

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type Extension struct {
	InteractionsArray []string
}

type ExtensionParams struct {
	MakerAsset      string
	MakerAssetData  string
	TakerAssetData  string
	GetMakingAmount string
	GetTakingAmount string
	Predicate       string
	Permit          string
	PreInteraction  string
	PostInteraction string
}

func NewExtension(params ExtensionParams) (Extension, error) {

	if params.Permit != "" {
		if params.MakerAsset == "" {
			return Extension{}, fmt.Errorf("when Permit is present, a maker asset must also be defined requires MakerAsset")
		}
	}

	if params.MakerAsset != "" {
		if params.Permit == "" {
			return Extension{}, fmt.Errorf("when MakerAsset is present, a maker asset must also be defined requires Permit")
		}
	}

	makerAssetData := params.MakerAssetData
	takerAssetData := params.TakerAssetData
	getMakingAmount := params.GetMakingAmount
	getTakingAmount := params.GetTakingAmount
	predicate := params.Predicate
	permit := params.MakerAsset + strings.TrimPrefix(params.Permit, "0x")
	preInteraction := params.PreInteraction
	postInteraction := params.PostInteraction

	interactions := []string{makerAssetData, takerAssetData, getMakingAmount, getTakingAmount, predicate, permit, preInteraction, postInteraction}

	return Extension{
		InteractionsArray: interactions,
	}, nil
}

func (i *Extension) Encode() string {
	interactionsConcatednated := i.getConcatenatedInteractions()
	if interactionsConcatednated == "" {
		return "0x"
	}

	extensionIndented, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic("Error marshaling extension")
	}
	fmt.Printf("Extension indented: %s\n", string(extensionIndented))

	offsetsBytes := i.getOffsets()
	paddedOffsetHex := fmt.Sprintf("%064x", offsetsBytes)
	return "0x" + paddedOffsetHex + interactionsConcatednated
}

func (i *Extension) getConcatenatedInteractions() string {
	var builder strings.Builder
	for _, interaction := range i.InteractionsArray {
		interaction = strings.TrimPrefix(interaction, "0x")
		builder.WriteString(interaction)
	}
	return builder.String()
}

func (i *Extension) getOffsets() *big.Int {
	var lengthMap []int
	for _, interaction := range i.InteractionsArray {
		lengthMap = append(lengthMap, len(strings.TrimPrefix(interaction, "0x"))/2)
	}

	cumulativeSum := 0
	bytesAccumulator := big.NewInt(0)
	var index uint64

	for _, length := range lengthMap {
		cumulativeSum += length
		shiftVal := big.NewInt(int64(cumulativeSum))
		shiftVal.Lsh(shiftVal, uint(32*index))           // Shift left
		bytesAccumulator.Add(bytesAccumulator, shiftVal) // Add to accumulator
		index++
	}

	return bytesAccumulator
}
