package http_executor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUserAgent(t *testing.T) {
	tests := []struct {
		name           string
		expectedPrefix string
	}{
		{
			name:           "User agent carries the client name and a version",
			expectedPrefix: "1inch-dev-portal-client-go:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := buildUserAgent()
			assert.True(t, strings.HasPrefix(result, tc.expectedPrefix), "user agent %q must start with %q", result, tc.expectedPrefix)
			assert.NotEqual(t, tc.expectedPrefix, result, "user agent must include a version suffix")
		})
	}
}
