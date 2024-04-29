package internal

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
			description:     "Success with v = 0x1b",
			signature:       "4ca2e082038ead998ff272272153b02643135a47bf8a820e5cebb5303763a3211117316de39fbe43781e78e60db312c9173ee7da13fec73ef8366af0d89c9f3c1b",
			expectedRValue:  "4ca2e082038ead998ff272272153b02643135a47bf8a820e5cebb5303763a321",
			expectedVSValue: "1117316de39fbe43781e78e60db312c9173ee7da13fec73ef8366af0d89c9f3c",
		},
		{
			description:     "Success with v = 0x1c",
			signature:       "2fac11bfe002d84bd0837f6efc88688bf4a35309bb5cfde80f740105ddbc9e024e552465e5087d9997739ba467e161c9752364d16cebaf9afd9f8e1a8f22addc1c",
			expectedRValue:  "2fac11bfe002d84bd0837f6efc88688bf4a35309bb5cfde80f740105ddbc9e02",
			expectedVSValue: "ce552465e5087d9997739ba467e161c9752364d16cebaf9afd9f8e1a8f22addc",
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
