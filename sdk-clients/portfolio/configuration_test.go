package portfolio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigurationAPI(t *testing.T) {
	configAPI, err := NewConfiguration(ConfigurationParams{
		ApiUrl: "https://api.example.com",
		ApiKey: "apikey123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, configAPI)
	assert.Equal(t, "https://api.example.com", configAPI.ApiURL)
	assert.Equal(t, "apikey123", configAPI.ApiKey)
}
