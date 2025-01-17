package keccak

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeccak256Legacy(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "Empty input",
			input: []byte(""),
			want:  "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
		{
			name:  "hello",
			input: []byte("hello"),
			want:  "0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8",
		},
		{
			name:  "abc",
			input: []byte("abc"),
			want:  "0x4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45",
		},
		{
			name:  "some random input",
			input: []byte("some random input"),
			want:  "0xd9092dd563673ef538f6d6162a5d127ddbd3d220a3f1f0eff6f7bac0194acd19",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Keccak256Legacy(tt.input)
			require.Equal(t, tt.want, got, "hash should match expected result")
		})
	}
}
