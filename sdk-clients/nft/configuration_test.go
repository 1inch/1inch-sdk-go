package nft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigurationAPI(t *testing.T) {
	configAPI, err := NewConfiguration("https://api.example.com", "apikey123")

	assert.NoError(t, err)
	assert.NotNil(t, configAPI)
	assert.Equal(t, "https://api.example.com", configAPI.ApiURL)
	assert.Equal(t, "apikey123", configAPI.ApiKey)
}
