package bigint

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

var largeString, _ = big.NewInt(0).SetString("9999999999999999999999999999", 10)

func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *big.Int
		wantErr bool
	}{
		{
			name:    "Valid decimal (10)",
			input:   "10",
			want:    big.NewInt(10),
			wantErr: false,
		},
		{
			name:    "Zero value",
			input:   "0",
			want:    big.NewInt(0),
			wantErr: false,
		},
		{
			name:    "Negative value (-100)",
			input:   "-100",
			want:    big.NewInt(-100),
			wantErr: false,
		},
		{
			name:    "Large number",
			input:   "9999999999999999999999999999",
			want:    largeString,
			wantErr: false,
		},
		{
			name:    "Invalid string",
			input:   "hello",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.input)
			if tt.wantErr {
				// We expect an error
				require.Error(t, err, "Expected an error but got none")
				return
			}

			// We do NOT expect an error
			require.NoError(t, err, "Did not expect an error but got one: %v", err)

			// Verify the returned *big.Int matches our expectation
			// (a nil `tt.want` is unusual when wantErr=false, so ensure it's not nil)
			require.NotNil(t, got, "Expected non-nil big.Int but got nil")

			// Compare numerical values
			if got.Cmp(tt.want) != 0 {
				t.Errorf("FromString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
