package http_executor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/1inch/1inch-sdk-go/v5/internal/version"
)

func TestUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "User agent carries the client name and the release version",
			expected: "1inch-dev-portal-client-go:" + version.Version,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, userAgent)
		})
	}
}
