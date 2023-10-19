package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
				ApiKey:            "",
			},
			expectedEnvironment:      baseUrlProduction.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Production (excluded entry)",
			config: Config{
				ApiKey: "",
			},
			expectedEnvironment:      baseUrlProduction.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Staging",
			config: Config{
				TargetEnvironment: EnvironmentStaging,
				ApiKey:            "",
			},
			expectedEnvironment:      baseUrlStaging.Host,
			expectedErrorDescription: "",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			client, err := NewClient(tc.config)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedEnvironment, client.BaseURL.Host)
		})
	}
}
