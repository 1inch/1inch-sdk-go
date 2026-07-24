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

func TestReleaseVersionRegex(t *testing.T) {
	tests := []struct {
		name    string
		version string
		allowed bool
	}{
		{
			name:    "Release tag",
			version: "v4.1.0",
			allowed: true,
		},
		{
			name:    "Multi-digit release tag",
			version: "v10.22.333",
			allowed: true,
		},
		{
			name:    "Pseudo-version",
			version: "v4.1.1-0.20260801120000-abcdef123456",
			allowed: false,
		},
		{
			name:    "Prerelease tag",
			version: "v4.2.0-rc.1",
			allowed: false,
		},
		{
			name:    "Toolchain placeholder",
			version: "(devel)",
			allowed: false,
		},
		{
			name:    "Build metadata suffix",
			version: "v4.1.0+incompatible",
			allowed: false,
		},
		{
			name:    "Empty version",
			version: "",
			allowed: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.allowed, releaseVersionRegex.MatchString(tc.version))
		})
	}
}
