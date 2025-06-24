package orderbook

import (
	"testing"

	web3_provider "github.com/1inch/1inch-sdk-go/internal/web3-provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSalt(t *testing.T) {
	tests := []struct {
		name         string
		extension    Extension
		expectedSalt string
	}{
		{
			name: "Salt",
			extension: Extension{
				MakerAssetSuffix: "0x",
				TakerAssetSuffix: "0x",
				MakingAmountData: "0x2ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859e6260000b401bf92000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
				TakingAmountData: "0x2ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859e6260000b401bf92000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
				Predicate:        "0x",
				MakerPermit:      "0x",
				PreInteraction:   "0x",
				PostInteraction:  "0x2ad5004c60e16e54d5007c80ce329adde5b51ef500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000646859e6150ac0866635457d36ab318d000000000000000000000000000066593d4e7d3a5f55167f0000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d0000972ad4499f120902631a000095770895ad27ad6b0d9500000000000000000000000000000000000000000000000000000000000000075dec5a",
			},
			expectedSalt: "0x2677cec45f20782506454895743b07ed0eae652cb39033bb6e4a3c7fa8662b5c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			encoded, err := tc.extension.Encode()
			require.NoError(t, err)

			salt, err := GenerateSalt(encoded)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedSalt[len(tc.expectedSalt)-40:], salt[len(salt)-40:])
		})
	}
}

func TestPrivateKeyProviderLive(t *testing.T) {
	testPrivateKey := "d8d1f95deb28949ea0ecc4e9a0decf89e98422c2d76ab6e5f736792a388c56c7"

	wallet, err := web3_provider.DefaultWalletOnlyProvider(testPrivateKey, 1)
	require.NoError(t, err)

	createOrderParams := CreateOrderParams{
		Wallet:             wallet,
		Salt:               "618054093254",
		MakerAsset:         "0xe9e7cea3dedca5984780bafc599bd69add087d56",
		TakerAsset:         "0x111111111117dc0aa78b770fa6a738034120c302",
		Maker:              "0xfb3c7eb936cAA12B5A884d612393969A557d4307",
		Taker:              "0x0000000000000000000000000000000000000000",
		MakingAmount:       "1000000000000000000",
		TakingAmount:       "1000000000000000000",
		MakerTraitsEncoded: "0",
		ExtensionEncoded:   "",
	}

	order, err := CreateLimitOrderMessage(createOrderParams, 1)
	require.NoError(t, err)

	expectedSignature := "0x8e1cbdc41ebb253aea91bfa41a028e735be4a5b25d93da0e3a6817070f40dcd31dfbc38bd3800ce2ff88089c77ca2f442dc84637006808aab0af00d966c917b11b"
	assert.Equal(t, expectedSignature, order.Signature)
}

func TestCreateLimitOrderMessage(t *testing.T) {

	//tests := []struct {
	//	name                string
	//	chainId             uint64
	//	makerTraitsAsString string
	//	createOrderParams   CreateOrderParams
	//	expectedOrderHash   string
	//	expectedSignature   string
	//}{
	//{
	//	name:                "Limit Order Creation",
	//	chainId:             137,
	//	makerTraitsAsString: "0x8a00000000000000000000005daf28d95c006859dd9f00000000000000000000",
	//	createOrderParams: CreateOrderParams{
	//		Wallet:      nil,
	//		SeriesNonce: nil,
	//		MakerTraits: &MakerTraits{
	//			AllowedSender:       "",
	//			Expiry:              1750719903,
	//			Nonce:               402370648412,
	//			Series:              0,
	//			NoPartialFills:      false,
	//			NeedPostinteraction: true,
	//			NeedPreinteraction:  false,
	//			NeedEpochCheck:      false,
	//			HasExtension:        true,
	//			ShouldUsePermit2:    false,
	//			ShouldUnwrapWeth:    false,
	//			AllowPartialFills:   false,
	//			AllowMultipleFills:  false,
	//		},
	//		Extension: Extension{
	//			MakerAssetSuffix: "0x",
	//			TakerAssetSuffix: "0x",
	//			MakingAmountData: "0xabd4e5fb590aa132749bbf2a04ea57efbaac399e000000000000006859dceb0000b401bf9f000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
	//			TakingAmountData: "0xabd4e5fb590aa132749bbf2a04ea57efbaac399e000000000000006859dceb0000b401bf9f000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
	//			Predicate:        "",
	//			MakerPermit:      "",
	//			PreInteraction:   "",
	//			PostInteraction:  "0xabd4e5fb590aa132749bbf2a04ea57efbaac399e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000646859dceb0ac0866635457d36ab318d000000000000000000000000000066593d4e7d3a5f55167f0000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d0000972ad4499f120902631a000095770895ad27ad6b0d95000000000000000000000000000000000000000000000000000000000000000752ba5a",
	//		},
	//		Maker:                          "0x50c5df26654B5EFBdD0c54a062dfa6012933deFe",
	//		MakerAsset:                     "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619",
	//		TakerAsset:                     "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359",
	//		TakingAmount:                   "469145",
	//		MakingAmount:                   "200000000000000",
	//		Taker:                          "0x0000000000000000000000000000000000000000",
	//		SkipWarnings:                   false,
	//		EnableOnchainApprovalsIfNeeded: false,
	//	},
	//	expectedSignature: "0xde99687c8e578779bbf9cdad4e6ffad4bc8e700e7f088242c0a983421458904b6f43dbaadf21f54bdd5dadbdd5aca8cedb920679c969b454eed948199bf0865d1b",
	//},
	//{
	//	name:                "Limit Order Creation JS SDK",
	//	chainId:             137,
	//	makerTraitsAsString: "0x8a0000000000000000000000775e2395d9006859e06300000000000000000000",
	//	createOrderParams: CreateOrderParams{
	//		Wallet:      nil,
	//		SeriesNonce: nil,
	//		MakerTraits: &MakerTraits{
	//			AllowedSender:       "",
	//			Expiry:              1750719903,
	//			Nonce:               402370648412,
	//			Series:              0,
	//			NoPartialFills:      false,
	//			NeedPostinteraction: true,
	//			NeedPreinteraction:  false,
	//			NeedEpochCheck:      false,
	//			HasExtension:        true,
	//			ShouldUsePermit2:    false,
	//			ShouldUnwrapWeth:    false,
	//			AllowPartialFills:   false,
	//			AllowMultipleFills:  false,
	//		},
	//		Extension: Extension{
	//			MakerAssetSuffix: "0x",
	//			TakerAssetSuffix: "0x",
	//			MakingAmountData: "0x2ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859dfa30000b401bf8e000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
	//			TakingAmountData: "0x2ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859dfa30000b401bf8e000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
	//			Predicate:        "",
	//			MakerPermit:      "",
	//			PreInteraction:   "",
	//			PostInteraction:  "0x2ad5004c60e16e54d5007c80ce329adde5b51ef500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000646859df920ac0866635457d36ab318d000000000000000000000000000066593d4e7d3a5f55167f0000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d0000972ad4499f120902631a000095770895ad27ad6b0d95000000000000000000000000000000000000000000000000000000000000000758295a",
	//		},
	//		Maker:                          "0x50c5df26654b5efbdd0c54a062dfa6012933defe",
	//		MakerAsset:                     "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
	//		TakerAsset:                     "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
	//		TakingAmount:                   "470525",
	//		MakingAmount:                   "200000000000000",
	//		Taker:                          "0x0000000000000000000000000000000000000000",
	//		SkipWarnings:                   false,
	//		EnableOnchainApprovalsIfNeeded: false,
	//	},
	//	expectedOrderHash: "0xd0e7bf0fe6d711f6f939c5262a17abfcbbe686b8cb44b08bdc8978aaa1e6d44d",
	//	expectedSignature: "0x8b82450612bd3552c389396e945553a52acc496a86c843742511bae67565519a4a5ceeeb5e0cb55568e9fc161c3e2e35a77842a6a0579e520ef9333a35e6b7881c",
	//},
	//}
	//
	//for _, tc := range tests {
	//	t.Run(tc.name, func(t *testing.T) {
	//
	//		oldMonkeyFunc := monkeyFunc
	//		monkeyFunc = func() string {
	//			return tc.makerTraitsAsString
	//		}
	//		defer func() {
	//			monkeyFunc = oldMonkeyFunc
	//		}()
	//
	//		w, err := web3_provider.DefaultWalletOnlyProvider(os.Getenv("WALLET_KEY"), tc.chainId)
	//		require.NoError(t, err)
	//
	//		tc.createOrderParams.Wallet = w
	//
	//		order, err := CreateLimitOrderMessage(tc.createOrderParams, int(tc.chainId))
	//		require.NoError(t, err)
	//		assert.Equal(t, tc.expectedSignature, order.Signature)
	//	})
	//}
}
