package nft

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetSupportedChains Returns supported chains
func (api *api) GetSupportedChains(ctx context.Context) (*SupportedChainsResponse, error) {
	u := fmt.Sprintf("/nft/v1/supportedchains")

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response SupportedChainsResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetNFTsByAddress Get users NFTs by EVM address with supported chains
func (api *api) GetNFTsByAddress(ctx context.Context, params GetNftsByAddressParams) (*GetNFTsByAddressResponse, error) {
	u := fmt.Sprintf("/nft/v1/byaddress")

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	stringChainIds := make([]string, len(params.ChainIds))

	for i, id := range params.ChainIds {
		stringChainIds[i] = strconv.Itoa(int(id))
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: struct {
			ChainIds string `url:"chainIds" json:"chainIds"`
			Address  string `url:"address" json:"address"`
			Limit    *int   `url:"limit,omitempty" json:"limit,omitempty"`
			Offset   *int   `url:"offset,omitempty" json:"offset,omitempty"`
		}{
			ChainIds: strings.Join(stringChainIds, ","),
			Address:  params.Address,
			Limit:    params.Limit,
			Offset:   params.Offset,
		},
		U:    u,
		Body: nil,
	}

	var response GetNFTsByAddressResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
