package multicall

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"
)

const multicallMethod = "multicall"

const (
	multicallContractEthereum  = "0x8d035edd8e09c3283463dade67cc0d49d6868063"
	multicallContractBnb       = "0x804708de7af615085203fa2b18eae59c5738e2a9"
	multicallContractPolygon   = "0x0196e8a9455a90d392b46df8560c867e7df40b34"
	multicallContractOptimism  = "0xE295aD71242373C37C5FdA7B57F26f9eA1088AFe"
	multicallContractArbitrum  = "0x11DEE30E710B8d4a8630392781Cc3c0046365d4c"
	multicallContractGnosis    = "0xE295aD71242373C37C5FdA7B57F26f9eA1088AFe"
	multicallContractAvalanche = "0xc4a8b7e29e3c8ec560cd4945c1cf3461a85a148d"
	multicallContractFantom    = "0xa31bb36c5164b165f9c36955ea4ccbab42b3b28e"
	multicallContractKlaytn    = "0xa31bb36c5164b165f9c36955ea4ccbab42b3b28e"
	multicallContractAurora    = "0xa0446d8804611944f1b527ecd37d7dcbe442caba"
	multicallContractZkSyncEra = "0xae1f66df155c611c15a23f31acf5a9bf1b87907e"
	multicallContractBase      = "0xa0446d8804611944f1b527ecd37d7dcbe442caba"
)

var (
	ErrEmptyContractAddr = errors.New("contract address must not be empty")
	ErrEmptyResponse     = errors.New("empty response")
)

type Client interface {
	ethereum.ContractCaller
	ethereum.ChainReader
}

func BuildCallData(to, data string, gas uint64, opts ...string) (r CallData) {
	r.To = to
	r.Data = data
	r.Gas = gas
	if len(opts) != 0 {
		r.MethodName = opts[0]
	}
	return r
}

func MultiCall(ctx context.Context, params MulticallParams) ([][]byte, error) {
	var requests []request
	for _, d := range params.Calldata {
		requests = append(requests, request{
			To:   common.HexToAddress(d.To),
			Data: common.FromHex(d.Data),
		})
	}

	multicallContract, err := abi.JSON(strings.NewReader(Multicallv2abiABI)) // Make a generic version of this ABI
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi error: %s", err)
	}

	data, err := multicallContract.Pack(
		multicallMethod,
		requests,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack msg error: %s", err)
	}

	evmHelperContract, err := GetEvmHelperContract(params.ChainId)
	if err != nil {
		return nil, err
	}

	var multicallResponse response
	resp, err := Call(ctx, CallParams{
		Client:          params.Client,
		Data:            data,
		ContractAddress: evmHelperContract,
		Block:           nil, // nil block means latest block
	})
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, ErrEmptyResponse
	}

	err = multicallContract.UnpackIntoInterface(&multicallResponse, multicallMethod, resp)
	if err != nil {
		return nil, err
	}

	return multicallResponse.Results, nil
}

func Call(ctx context.Context, params CallParams) ([]byte, error) {
	if params.ContractAddress == "" {
		return nil, ErrEmptyContractAddr
	}
	var toAddress = common.HexToAddress(params.ContractAddress)
	var nodeMsg = ethereum.CallMsg{
		To:   &toAddress,
		Data: params.Data,
	}
	resp, err := params.Client.CallContract(ctx, nodeMsg, params.Block)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: error: %s", err)
	}

	return resp, nil
}

func GetEvmHelperContract(chainId int) (string, error) {
	switch chainId {
	case chains.Ethereum:
		return multicallContractEthereum, nil
	case chains.Bsc:
		return multicallContractBnb, nil
	case chains.Polygon:
		return multicallContractPolygon, nil
	case chains.Optimism:
		return multicallContractOptimism, nil
	case chains.Arbitrum:
		return multicallContractArbitrum, nil
	case chains.Gnosis:
		return multicallContractGnosis, nil
	case chains.Avalanche:
		return multicallContractAvalanche, nil
	case chains.Fantom:
		return multicallContractFantom, nil
	case chains.Klaytn:
		return multicallContractKlaytn, nil
	case chains.Aurora:
		return multicallContractAurora, nil
	case chains.ZkSyncEra:
		return multicallContractZkSyncEra, nil
	case chains.Base:
		return multicallContractBase, nil
	default:
		return "", fmt.Errorf("chain %d is not supported", chainId)
	}
}
