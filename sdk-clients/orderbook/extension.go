package orderbook

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/1inch/1inch-sdk-go/internal/bytesiterator"
	"github.com/ethereum/go-ethereum/common/math"
)

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

func NewExtension(params ExtensionParams) (*Extension, error) {

	if params.Permit != "" {
		if params.MakerAsset == "" {
			return nil, fmt.Errorf("when Permit is present, a maker asset must also be defined requires MakerAsset")
		}
	}

	if params.MakerAsset != "" {
		if params.Permit == "" {
			return nil, fmt.Errorf("when MakerAsset is present, a maker asset must also be defined requires Permit")
		}
	}

	return &Extension{
		MakerAssetSuffix: params.MakerAssetData,
		TakerAssetSuffix: params.TakerAssetData,
		MakingAmountData: params.GetMakingAmount,
		TakingAmountData: params.GetTakingAmount,
		Predicate:        params.Predicate,
		MakerPermit:      params.MakerAsset + strings.TrimPrefix(params.Permit, "0x"),
		PreInteraction:   params.PreInteraction,
		PostInteraction:  params.PostInteraction,
	}, nil
}

type Extension struct {
	MakerAssetSuffix string
	TakerAssetSuffix string
	MakingAmountData string
	TakingAmountData string
	Predicate        string
	MakerPermit      string
	PreInteraction   string
	PostInteraction  string
}

// Decode decodes the input byte slice into an Extension struct using reflection.
func Decode(data []byte) (*Extension, error) {
	// TODO Handle the special case where data equals ZX.
	//if string(data) == ZX {
	//	return DefaultExtension(), nil
	//}

	iter := bytesiterator.New(data)

	// Read the first 32 bytes as offsets.
	offsets, err := iter.NextUint256()
	if err != nil {
		return &Extension{}, errors.New("failed to read offsets: " + err.Error())
	}

	consumed := 0

	// Initialize the Extension struct
	var ext Extension

	// Use reflection to iterate over the struct fields in order.
	val := reflect.ValueOf(&ext).Elem() // Get the reflect.Value of the struct
	typ := val.Type()                   // Get the reflect.Type of the struct

	numFields := typ.NumField()

	// Iterate through all fields except the last one (CustomData)
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip CustomData for now
		if field.Name == "CustomData" {
			continue
		}

		const uint32Max = math.MaxUint32

		// Extract the lowest 32 bits for the current field's offset.
		offset := new(big.Int).And(offsets, big.NewInt(uint32Max)).Uint64()
		bytesCount := int(offset) - consumed

		if bytesCount < 0 {
			return &Extension{}, errors.New("invalid offset leading to negative bytesCount for field: " + field.Name)
		}

		// Read the next bytesCount bytes for the current field.
		fieldBytes, err := iter.NextBytes(bytesCount)
		if err != nil {
			return &Extension{}, errors.New("failed to read field " + field.Name + ": " + err.Error())
		}
		if len(fieldBytes) < bytesCount {
			return &Extension{}, errors.New("insufficient bytes for field " + field.Name)
		}

		// Set the field value using reflection.
		if field.Type.Kind() == reflect.String {
			fieldVal.SetString(fmt.Sprintf("%x", fieldBytes))
		} else {
			return &Extension{}, errors.New("unsupported field type for field: " + field.Name)
		}

		// Update the consumed bytes and shift the offsets for the next field.
		consumed += bytesCount
		offsets = new(big.Int).Rsh(offsets, 32)
	}

	// TODO The remaining bytes are considered as CustomData, but it is not supported yet.
	//customDataBytes, err := iter.Rest()
	//if err != nil {
	//	return &Extension{}, errors.New("failed to read CustomData: " + err.Error())
	//}
	//ext.CustomData = string(customDataBytes)

	return &ext, nil
}

// hexToBytes converts a hexadecimal string to a byte slice.
func hexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// contains checks if the substring is present in the string.
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

// Encode encodes the Extension struct into a hex string with offsets.
func (ext *Extension) Encode() (string, error) {
	fields := []string{
		ext.MakerAssetSuffix,
		ext.TakerAssetSuffix,
		ext.MakingAmountData,
		ext.TakingAmountData,
		ext.Predicate,
		ext.MakerPermit,
		ext.PreInteraction,
		ext.PostInteraction,
	}

	var byteCounts []int
	var dataBytes []byte

	// Decode each field and collect byte counts
	for _, field := range fields {
		fieldStr := strings.TrimPrefix(field, "0x")
		fieldStr = strings.TrimPrefix(fieldStr, "0X")

		// Ensure even length for hex decoding
		if len(fieldStr)%2 != 0 {
			fieldStr = "0" + fieldStr
		}

		fieldData, err := hex.DecodeString(fieldStr)
		if err != nil {
			return "", fmt.Errorf("failed to decode field '%s': %v", field, err)
		}
		byteCounts = append(byteCounts, len(fieldData))
		dataBytes = append(dataBytes, fieldData...)
	}

	// Calculate cumulative offsets
	cumulativeSum := 0
	var offsets []byte
	for i := 0; i < len(byteCounts); i++ {
		cumulativeSum += byteCounts[i]
		offsetBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(offsetBytes, uint32(cumulativeSum))
		offsets = append(offsetBytes, offsets...)
	}

	// Encode offsets and data to hex
	offsetsHex := hex.EncodeToString(offsets)
	dataHex := hex.EncodeToString(dataBytes)

	// Concatenate with "0x" prefix
	return "0x" + offsetsHex + dataHex, nil
}
