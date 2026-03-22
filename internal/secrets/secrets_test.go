package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupSecrets_ValidKey(t *testing.T) {
	apiKey = ""
	t.Setenv("API_KEY", "test-api-key-123")

	err := SetupSecrets()
	require.NoError(t, err)

	key, err := GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, "test-api-key-123", key)
}

func TestSetupSecrets_MissingKey(t *testing.T) {
	apiKey = ""
	t.Setenv("API_KEY", "")

	err := SetupSecrets()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API_KEY")
}

func TestGetAPIKey_BeforeSetup(t *testing.T) {
	apiKey = ""

	key, err := GetAPIKey()
	assert.Error(t, err)
	assert.Empty(t, key)
	assert.Contains(t, err.Error(), "not initialized")
}
