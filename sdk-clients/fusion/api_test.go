package fusion

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
)

type MockHttpExecutor struct {
	Called      bool
	ExecuteErr  error
	ResponseObj interface{}
}

func (m *MockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v interface{}) error {
	m.Called = true
	if m.ExecuteErr != nil {
		return m.ExecuteErr
	}

	// Copy the mock response object to v
	if m.ResponseObj != nil && v != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("v must be a non-nil pointer")
		}
		reflect.Indirect(rv).Set(reflect.ValueOf(m.ResponseObj))
	}
	return nil
}

func TestGetNFTsByAddress(t *testing.T) {
	ctx := context.Background()

	mockedResp := GetNFTsByAddressResponse{
		Assets: []AssetExtended{
			{
				ID:                   18437292988043788,
				Provider:             "POAP",
				TokenID:              "6925855",
				AnimationOriginalURL: "https://app.poap.xyz/token/6925855",
				AnimationURL:         "https://assets.poap.xyz/297f83c9-47e9-46d2-84fc-215ec2da8bb4.gif",
				Description:          "Join our crazy CryptoCanal community of Amsterdam at Two Chefs foodbar for a couple of beers and chats. Old friends and new faces welcome! 🦦\n\n****************************\nMeetup sponsors\nBitvavo, the biggest crypto exchange in the Netherlands!\nVanEck, founded in 1955, they offer innovative passive and active investment strategies.\n\n****************************\nCryptoCanal offer education, event and consultancy services for the crypto industry and we're not afraid to be political. We organise ETHDam, an annual conference and hackathon in the heart Amsterdam.\nFollow us on Twitter and join our Telegram group for the daily updates.\nNB. Pictures and videos that might be used for promotional purposes will be taken during the event.",
				ExternalLink:         "https://app.poap.xyz/token/6925855",
				Permalink:            "https://app.poap.xyz/token/6925855",
				Name:                 "Crypto Drinks Amsterdam by CryptoCanal",
				ChainID:              100,
				Traits: []Trait{
					{TraitType: "startDate", Value: "28-Nov-2023"},
					{TraitType: "endDate", Value: "29-Nov-2023"},
					{TraitType: "virtualEvent", Value: "false"},
					{TraitType: "city", Value: "Amsterdam"},
					{TraitType: "country", Value: "Netherlands"},
					{TraitType: "eventURL", Value: "https://lu.ma/cryptodrinksamsnov"},
				},
				Priority: 2,
				AssetContract: AssetContractExtended{
					Address:           "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
					AssetContractType: "",
					CreatedDate:       "2023-11-29 19:04:40",
					Name:              "Crypto Drinks Amsterdam by CryptoCanal",
					Description:       "Join our crazy CryptoCanal community of Amsterdam at Two Chefs foodbar for a couple of beers and chats. Old friends and new faces welcome! 🦦\n\n****************************\nMeetup sponsors\nBitvavo, the biggest crypto exchange in the Netherlands!\nVanEck, founded in 1955, they offer innovative passive and active investment strategies.\n\n****************************\nCryptoCanal offer education, event and consultancy services for the crypto industry and we're not afraid to be political. We organise ETHDam, an annual conference and hackathon in the heart Amsterdam.\nFollow us on Twitter and join our Telegram group for the daily updates.\nNB. Pictures and videos that might be used for promotional purposes will be taken during the event.",
					ImageURL:          "https://assets.poap.xyz/297f83c9-47e9-46d2-84fc-215ec2da8bb4.gif",
					ExternalLink:      "https://api.poap.tech/metadata/159773/6925855",
					TotalSupply:       "12",
					Owner:             "159773",
					OpenseaVersion:    "",
					NFTVersion:        "",
					SchemaName:        "ERC721",
					Symbol:            "",
				},
			},
			{
				ID:                   209013196794191680,
				Provider:             "POAP",
				TokenID:              "6717010",
				AnimationOriginalURL: "https://app.poap.xyz/token/6717010",
				AnimationURL:         "https://assets.poap.xyz/ethcc-5b65d-attendee-2023-logo-1688660565554.gif",
				Description:          "You've attended the 6th edition of the Ethereum Community Conference, held in Paris at the Maison de la Mutuliaté and the Collège des Bernardins from 17 to 20 July 2023. Thank you for your participation.\n\nConnect with other Attendees in the SALSA poap-gated chatroom - https://app.salsa.me/ethcc",
				ExternalLink:         "https://app.poap.xyz/token/6717010",
				Permalink:            "https://app.poap.xyz/token/6717010",
				Name:                 "EthCC[6] - Attendee",
				ChainID:              100,
				Traits: []Trait{
					{TraitType: "startDate", Value: "17-Jul-2023"},
					{TraitType: "endDate", Value: "20-Jul-2023"},
					{TraitType: "virtualEvent", Value: "false"},
					{TraitType: "city", Value: "Paris"},
					{TraitType: "country", Value: "France"},
					{TraitType: "eventURL", Value: "https://app.salsa.me/ethcc"},
				},
				Priority: 2,
				AssetContract: AssetContractExtended{
					Address:           "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
					AssetContractType: "",
					CreatedDate:       "2023-07-17 11:26:10",
					Name:              "EthCC[6] - Attendee",
					Description:       "You've attended the 6th edition of the Ethereum Community Conference, held in Paris at the Maison de la Mutuliaté and the Collège des Bernardins from 17 to 20 July 2023. Thank you for your participation.\n\nConnect with other Attendees in the SALSA poap-gated chatroom - https://app.salsa.me/ethcc",
					ImageURL:          "https://assets.poap.xyz/ethcc-5b65d-attendee-2023-logo-1688660565554.gif",
					ExternalLink:      "https://api.poap.tech/metadata/141910/6717010",
					TotalSupply:       "1297",
					Owner:             "141910",
					OpenseaVersion:    "",
					NFTVersion:        "",
					SchemaName:        "ERC721",
					Symbol:            "",
				},
			},
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
	}

	params := GetNftsByAddressParams{
		ChainIds: []GetNftsByAddressParamsChainIds{
			constants.EthereumChainId,
			constants.GnosisChainId,
		},
		Address: "0x083fc10cE7e97CaFBaE0fE332a9c4384c5f54E45",
	}

	prices, err := api.GetNFTsByAddress(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, prices)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(prices, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", prices, mockedResp)
	}
	require.NoError(t, err)
}

func TestGetSupportedChains(t *testing.T) {
	ctx := context.Background()

	mockedResp := SupportedChainsResponse([]GetNftsByAddressParamsChainIds{
		constants.EthereumChainId,
		constants.GnosisChainId,
	})

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
	}

	prices, err := api.GetSupportedChains(ctx)
	require.NoError(t, err)
	require.NotNil(t, prices)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(prices, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", prices, mockedResp)
	}
	require.NoError(t, err)
}
