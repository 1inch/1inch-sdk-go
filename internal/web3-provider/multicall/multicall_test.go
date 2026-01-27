package multicall

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestBuildCallData(t *testing.T) {
	tests := []struct {
		name     string
		to       common.Address
		data     []byte
		gas      uint64
		opts     []string
		expected CallData
	}{
		{
			name: "basic data",
			to:   common.HexToAddress("0xAb1234cdE56789f0Ab1234cdE56789f0ab1234cD"),
			data: []byte{0xde, 0xad, 0xbe, 0xef},
			gas:  21000,
			opts: nil,
			expected: CallData{
				To:   "0xAb1234cdE56789F0AB1234cDe56789F0AB1234CD",
				Data: "0xdeadbeef",
				Gas:  21000,
			},
		},
		{
			name: "with method name",
			to:   common.HexToAddress("0xAb1234cdE56789f0Ab1234cdE56789f0ab1234cD"),
			data: []byte{0xca, 0xfe, 0xba, 0xbe},
			gas:  30000,
			opts: []string{"transfer"},
			expected: CallData{
				To:         "0xAb1234cdE56789F0AB1234cDe56789F0AB1234CD",
				Data:       "0xcafebabe",
				Gas:        30000,
				MethodName: "transfer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := BuildCallData(tt.to, tt.data, tt.gas, tt.opts...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
