package multicall

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk-go/constants"
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
	ErrEmptyResponse = errors.New("empty response")
)

type Client interface {
	ethereum.ContractCaller
	ethereum.ChainReader
}

type Multicall struct {
	client          *ethclient.Client
	contractAddress *common.Address
	contractABI     *abi.ABI
}

func NewMulticall(client *ethclient.Client, chainId uint64) (*Multicall, error) {
	var addressRaw string

	switch chainId {
	case constants.EthereumChainId:
		addressRaw = multicallContractEthereum
	case constants.BscChainId:
		addressRaw = multicallContractBnb
	case constants.PolygonChainId:
		addressRaw = multicallContractPolygon
	case constants.OptimismChainId:
		addressRaw = multicallContractOptimism
	case constants.ArbitrumChainId:
		addressRaw = multicallContractArbitrum
	case constants.GnosisChainId:
		addressRaw = multicallContractGnosis
	case constants.AvalancheChainId:
		addressRaw = multicallContractAvalanche
	case constants.FantomChainId:
		addressRaw = multicallContractFantom
	case constants.KlaytnChainId:
		addressRaw = multicallContractKlaytn
	case constants.AuroraChainId:
		addressRaw = multicallContractAurora
	case constants.ZkSyncEraChainId:
		addressRaw = multicallContractZkSyncEra
	case constants.BaseChainId:
		addressRaw = multicallContractBase
	default:
		return nil, fmt.Errorf("chain %d is not supported", chainId)
	}

	helperContractAddress := common.HexToAddress(addressRaw)
	contractABI, err := abi.JSON(strings.NewReader(Multicallv2abiABI)) // Make a generic version of this ABI
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi error: %s", err)
	}
	return &Multicall{
		client:          client,
		contractAddress: &helperContractAddress,
		contractABI:     &contractABI,
	}, nil
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

func (m Multicall) Execute(ctx context.Context, callData []CallData) ([][]byte, error) {
	var requests []request
	for _, d := range callData {
		requests = append(requests, request{
			To:   common.HexToAddress(d.To),
			Data: common.FromHex(d.Data),
		})
	}

	data, err := m.contractABI.Pack(
		multicallMethod,
		requests,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack msg error: %s", err)
	}

	nodeMsg := ethereum.CallMsg{
		To:   m.contractAddress,
		Data: data,
	}
	resp, err := m.client.CallContract(ctx, nodeMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: error: %s", err)
	}

	if len(resp) == 0 {
		return nil, ErrEmptyResponse
	}

	var multicallResponse response
	err = m.contractABI.UnpackIntoInterface(&multicallResponse, multicallMethod, resp)
	if err != nil {
		return nil, err
	}

	return multicallResponse.Results, nil
}
