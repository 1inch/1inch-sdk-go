package orderbook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOrderIntegration(t *testing.T) {

	testcases := []struct {
		description     string
		signature       string
		expectedRValue  string
		expectedVSValue string
	}{
		{
			description:     "Success",
			signature:       "4ca2e082038ead998ff272272153b02643135a47bf8a820e5cebb5303763a3211117316de39fbe43781e78e60db312c9173ee7da13fec73ef8366af0d89c9f3c1b",
			expectedRValue:  "4ca2e082038ead998ff272272153b02643135a47bf8a820e5cebb5303763a321",
			expectedVSValue: "1117316de39fbe43781e78e60db312c9173ee7da13fec73ef8366af0d89c9f3c",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {

			compactSignature, err := CompressSignature(tc.signature)
			require.NoError(t, err)

			require.Equal(t, tc.expectedRValue, fmt.Sprintf("%x", compactSignature.R))
			require.Equal(t, tc.expectedVSValue, fmt.Sprintf("%x", compactSignature.VS))
		})
	}
}
