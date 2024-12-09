package fusionplus

//
//import (
//	"math/big"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//
//	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
//)
//
//func TestGenerateSalt(t *testing.T) {
//	// Save the original function
//	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
//
//	// Monkey patch the function
//	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
//		return big.NewInt(123456), nil
//	}
//
//	// Restore the original function after the test
//	defer func() {
//		random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc
//	}()
//
//	tests := []struct {
//		name      string
//		extension *Extension
//		expected  string
//		expectErr bool
//	}{
//		{
//			name: "Generate salt when extension is not empty",
//			extension: &Extension{
//				MakerAssetSuffix: "suffix1",
//				TakerAssetSuffix: "suffix2",
//				MakingAmountData: "data1",
//				TakingAmountData: "data2",
//				Predicate:        "predicate",
//				MakerPermit:      "permit",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//				CustomData:       "custom",
//			},
//			expected:  "180431658011416401710119735245975317914670388782711199",
//			expectErr: false,
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			expected, err := BigIntFromString(tc.expected)
//			require.NoError(t, err)
//
//			result, err := tc.extension.GenerateSalt()
//			if tc.expectErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//				assert.Equal(t, expected, result)
//			}
//		})
//	}
//}
//
//func TestNewExtension(t *testing.T) {
//	tests := []struct {
//		name      string
//		params    ExtensionParams
//		expectErr bool
//		errMsg    string
//	}{
//		{
//			name: "Valid parameters",
//			params: ExtensionParams{
//				MakerAssetSuffix: "0x1234",
//				TakerAssetSuffix: "0x1234",
//				MakingAmountData: "0x1234",
//				TakingAmountData: "0x1234",
//				Predicate:        "0x1234",
//				MakerPermit:      "0x1234",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//			},
//			expectErr: false,
//		},
//		{
//			name: "Invalid MakerAssetSuffix",
//			params: ExtensionParams{
//				MakerAssetSuffix: "invalid",
//				TakerAssetSuffix: "0x1234",
//				MakingAmountData: "0x1234",
//				TakingAmountData: "0x1234",
//				Predicate:        "0x1234",
//				MakerPermit:      "0x1234",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//			},
//			expectErr: true,
//			errMsg:    "MakerAssetSuffix must be valid hex string",
//		},
//		{
//			name: "Invalid TakerAssetSuffix",
//			params: ExtensionParams{
//				MakerAssetSuffix: "0x1234",
//				TakerAssetSuffix: "invalid",
//				MakingAmountData: "0x1234",
//				TakingAmountData: "0x1234",
//				Predicate:        "0x1234",
//				MakerPermit:      "0x1234",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//			},
//			expectErr: true,
//			errMsg:    "TakerAssetSuffix must be valid hex string",
//		},
//		{
//			name: "Invalid MakingAmountData",
//			params: ExtensionParams{
//				MakerAssetSuffix: "0x1234",
//				TakerAssetSuffix: "0x1234",
//				MakingAmountData: "invalid",
//				TakingAmountData: "0x1234",
//				Predicate:        "0x1234",
//				MakerPermit:      "0x1234",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//			},
//			expectErr: true,
//			errMsg:    "MakingAmountData must be valid hex string",
//		},
//		{
//			name: "Invalid TakingAmountData",
//			params: ExtensionParams{
//				MakerAssetSuffix: "0x1234",
//				TakerAssetSuffix: "0x1234",
//				MakingAmountData: "0x1234",
//				TakingAmountData: "invalid",
//				Predicate:        "0x1234",
//				MakerPermit:      "0x1234",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//			},
//			expectErr: true,
//			errMsg:    "TakingAmountData must be valid hex string",
//		},
//		{
//			name: "Invalid Predicate",
//			params: ExtensionParams{
//				MakerAssetSuffix: "0x1234",
//				TakerAssetSuffix: "0x1234",
//				MakingAmountData: "0x1234",
//				TakingAmountData: "0x1234",
//				Predicate:        "invalid",
//				MakerPermit:      "0x1234",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//			},
//			expectErr: true,
//			errMsg:    "Predicate must be valid hex string",
//		},
//		{
//			name: "Invalid MakerPermit",
//			params: ExtensionParams{
//				MakerAssetSuffix: "0x1234",
//				TakerAssetSuffix: "0x1234",
//				MakingAmountData: "0x1234",
//				TakingAmountData: "0x1234",
//				Predicate:        "0x1234",
//				MakerPermit:      "invalid",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//			},
//			expectErr: true,
//			errMsg:    "MakerPermit must be valid hex string",
//		},
//		{
//			name: "CustomData not supported",
//			params: ExtensionParams{
//				MakerAssetSuffix: "0x1234",
//				TakerAssetSuffix: "0x1234",
//				MakingAmountData: "0x1234",
//				TakingAmountData: "0x1234",
//				Predicate:        "0x1234",
//				MakerPermit:      "0x1234",
//				PreInteraction:   "pre",
//				PostInteraction:  "post",
//				CustomData:       "0x1234",
//			},
//			expectErr: true,
//			errMsg:    "CustomData is not currently supported",
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			ext, err := NewExtension(tc.params)
//			if tc.expectErr {
//				require.Error(t, err)
//				assert.Equal(t, tc.errMsg, err.Error())
//			} else {
//				require.NoError(t, err)
//				assert.NotNil(t, ext)
//				assert.Equal(t, tc.params.MakerAssetSuffix, ext.MakerAssetSuffix)
//				assert.Equal(t, tc.params.TakerAssetSuffix, ext.TakerAssetSuffix)
//				assert.Equal(t, tc.params.MakingAmountData, ext.MakingAmountData)
//				assert.Equal(t, tc.params.TakingAmountData, ext.TakingAmountData)
//				assert.Equal(t, tc.params.Predicate, ext.Predicate)
//				assert.Equal(t, tc.params.MakerPermit, ext.MakerPermit)
//				assert.Equal(t, tc.params.PreInteraction, ext.PreInteraction)
//				assert.Equal(t, tc.params.PostInteraction, ext.PostInteraction)
//				assert.Equal(t, tc.params.CustomData, ext.CustomData)
//			}
//		})
//	}
//}
