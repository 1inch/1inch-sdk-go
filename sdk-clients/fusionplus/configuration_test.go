package fusionplus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigurationAPI(t *testing.T) {
	configAPI, err := NewConfiguration(ConfigurationParams{
		ApiUrl:     "https://api.example.com",
		ApiKey:     "apikey123",
		PrivateKey: "965e092fdfc08940d2bd05c7b5c7e1c51e283e92c7f52bbf1408973ae9a9acb7",
	})

	require.NoError(t, err)
	assert.NotNil(t, configAPI)
	assert.Equal(t, "https://api.example.com", configAPI.APIConfiguration.ApiURL)
	assert.Equal(t, "apikey123", configAPI.APIConfiguration.ApiKey)
}
