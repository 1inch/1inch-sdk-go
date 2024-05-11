package gasprices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigurationAPI(t *testing.T) {
	configAPI, err := NewConfiguration(1, "https://api.example.com", "apikey123")

	assert.NoError(t, err)
	assert.NotNil(t, configAPI)
	assert.Equal(t, uint64(1), configAPI.API.chainId)
	assert.Equal(t, "https://api.example.com", configAPI.ApiURL)
	assert.Equal(t, "apikey123", configAPI.ApiKey)
}
