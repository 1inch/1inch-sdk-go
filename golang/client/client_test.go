package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		description              string
		config                   Config
		expectedEnvironment      string
		expectedErrorDescription string
	}{
		{
			description: "Production",
			config: Config{
				TargetEnvironment: EnvironmentProduction,
				ApiKey:            "abc123",
			},
			expectedEnvironment:      baseUrlProduction.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Production (excluded entry)",
			config: Config{
				ApiKey: "abc123",
			},
			expectedEnvironment:      baseUrlProduction.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Staging",
			config: Config{
				TargetEnvironment: EnvironmentStaging,
				ApiKey:            "abc123",
			},
			expectedEnvironment:      baseUrlStaging.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Error - unrecognized environment",
			config: Config{
				TargetEnvironment: Environment("invalid"),
				ApiKey:            "abc123",
			},
			expectedEnvironment:      baseUrlStaging.Host,
			expectedErrorDescription: "unrecognized environment: invalid",
		},
		{
			description: "Error - no API key",
			config: Config{
				ApiKey: "",
			},
			expectedEnvironment:      baseUrlProduction.Host,
			expectedErrorDescription: "API key is required",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, err := NewClient(tc.config)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedEnvironment, c.BaseURL.Host)
		})
	}
}
