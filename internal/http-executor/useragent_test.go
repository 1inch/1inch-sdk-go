package http_executor

import (
	"runtime/debug"
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

func TestSdkVersionFromBuildInfo(t *testing.T) {
	tests := []struct {
		name     string
		info     *debug.BuildInfo
		expected string
	}{
		{
			name: "Release tag from a consumer's dependency list",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/consumer", Version: "(devel)"},
				Deps: []*debug.Module{{Path: modulePath, Version: "v4.1.0"}},
			},
			expected: "v4.1.0",
		},
		{
			name: "Pseudo-version is reported verbatim",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/consumer"},
				Deps: []*debug.Module{{Path: modulePath, Version: "v4.1.1-0.20260801120000-abcdef123456"}},
			},
			expected: "v4.1.1-0.20260801120000-abcdef123456",
		},
		{
			name: "Prerelease tag is reported verbatim",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/consumer"},
				Deps: []*debug.Module{{Path: modulePath, Version: "v4.2.0-rc.1"}},
			},
			expected: "v4.2.0-rc.1",
		},
		{
			name: "Build metadata is reported verbatim",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/consumer"},
				Deps: []*debug.Module{{Path: modulePath, Version: "v4.1.0+incompatible"}},
			},
			expected: "v4.1.0+incompatible",
		},
		{
			name: "Versioned fork replacement reports the replacement version",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/consumer"},
				Deps: []*debug.Module{{
					Path:    modulePath,
					Version: "v4.1.0",
					Replace: &debug.Module{Path: "example.com/fork/v4", Version: "v4.1.2"},
				}},
			},
			expected: "v4.1.2",
		},
		{
			name: "Local path replacement reports unknown",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/consumer"},
				Deps: []*debug.Module{{
					Path:    modulePath,
					Version: "v4.1.0",
					Replace: &debug.Module{Path: "../1inch-sdk-go", Version: ""},
				}},
			},
			expected: "unknown",
		},
		{
			name: "In-repo build with a stamped release tag",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: modulePath, Version: "v4.1.0"},
			},
			expected: "v4.1.0",
		},
		{
			name: "In-repo build with the devel placeholder reports unknown",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: modulePath, Version: "(devel)"},
			},
			expected: "unknown",
		},
		{
			name:     "No build info reports unknown",
			info:     nil,
			expected: "unknown",
		},
		{
			name: "Unrelated dependencies report unknown",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/consumer"},
				Deps: []*debug.Module{{Path: "example.com/other", Version: "v1.0.0"}},
			},
			expected: "unknown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, sdkVersionFromBuildInfo(tc.info))
		})
	}
}
