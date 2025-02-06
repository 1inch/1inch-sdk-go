package aggregation

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"

	"github.com/1inch/1inch-sdk-go/constants"
)

func TestNormalizeSwapResponse(t *testing.T) {
	d := "0x0502b1c50000000000000000000000005a98fcbea516cf06857215779fd812ca3bef1b32000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000100000000000000003b6d0340c558f600b34a5f69dd2f0d06cb8a88d829b7420ade8bb62d"
	wantedDataLDOWETH, err := hexutil.Decode(d)
	assert.NoError(t, err)

	testCases := []struct {
		name    string
		input   SwapResponse
		want    *SwapResponseExtended
		wantErr bool
	}{
		{
			name: "LDO -> WETH ETH (small amount)",
			input: SwapResponse{
				Tx: TransactionData{
					Data:     d,
					From:     "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
					Gas:      257615,
					GasPrice: "22931145666",
					To:       "0x1111111254eeb25477b68fb85ed929f73a960582",
					Value:    "1000000000000000000",
				},
			},
			want: &SwapResponseExtended{
				SwapResponse: SwapResponse{
					Tx: TransactionData{
						Data:     d,
						From:     "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
						Gas:      257615,
						GasPrice: "22931145666",
						To:       "0x1111111254eeb25477b68fb85ed929f73a960582",
						Value:    "1000000000000000000",
					},
				},
				TxNormalized: NormalizedTransactionData{
					Data:     wantedDataLDOWETH,
					Gas:      257615,
					GasPrice: big.NewInt(22931145666),
					To:       common.HexToAddress("0x1111111254eeb25477b68fb85ed929f73a960582"),
					Value:    big.NewInt(1000000000000000000),
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid 'To' address",
			input: SwapResponse{
				Tx: TransactionData{
					Data:     "0x0502b1c50000000000000000000000005a98fcbea516cf06857215779fd812ca3bef1b32000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000100000000000000003b6d0340c558f600b34a5f69dd2f0d06cb8a88d829b7420ade8bb62d",
					From:     "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
					Gas:      257615,
					GasPrice: "22931145666",
					To:       "0xInvalid",
					Value:    "1000000000000000000",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid 'GasPrice'",
			input: SwapResponse{
				Tx: TransactionData{
					Data:     "0xdeadbeef",
					From:     "0x000000000000000000000000000000000000dead",
					Gas:      21000,
					GasPrice: "invalid",
					To:       "0x000000000000000000000000000000000000beef",
					Value:    "1000000000000000000",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeSwapResponse(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestNormalizeApproveCallDataResponse(t *testing.T) {
	d := "0x095ea7b30000000000000000000000001111111254eeb25477b68fb85ed929f73a9605820000000000000000000000000000000000000000000000000000000556cd83c2"
	ldoApproveData, err := hexutil.Decode(d)
	assert.NoError(t, err)

	testCases := []struct {
		name    string
		input   ApproveCallDataResponse
		want    *ApproveCallDataResponseExtended
		wantErr bool
	}{
		{
			name: "LDO Approve to 22931145666",
			input: ApproveCallDataResponse{
				Data:     d,
				GasPrice: "16955435273",
				To:       "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
				Value:    "0",
			},
			want: &ApproveCallDataResponseExtended{
				ApproveCallDataResponse: ApproveCallDataResponse{
					Data:     d,
					GasPrice: "16955435273",
					To:       "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
					Value:    "0",
				},
				TxNormalized: NormalizedTransactionData{
					Data:     ldoApproveData,
					Gas:      constants.ERC20_APPROVE_GAS,
					GasPrice: big.NewInt(16955435273),
					To:       common.HexToAddress("0x5a98fcbea516cf06857215779fd812ca3bef1b32"),
					Value:    big.NewInt(0),
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid 'To' Address",
			input: ApproveCallDataResponse{
				Data:     "0x095ea7b3",
				GasPrice: "100000000000",
				To:       "0xInvalidAddress",
				Value:    "0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid 'GasPrice'",
			input: ApproveCallDataResponse{
				Data:     "0x095ea7b3",
				GasPrice: "NotANumber",
				To:       "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
				Value:    "0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid 'Data'",
			input: ApproveCallDataResponse{
				Data:     "0xZZZ",
				GasPrice: "100000000000",
				To:       "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
				Value:    "0",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeApproveCallDataResponse(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
