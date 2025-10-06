package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretsHappyPath(t *testing.T) {
	apiKey = "" // Reset global state
	t.Setenv("API_KEY", "test-key")

	err := SetupSecrets()
	assert.NoError(t, err)

	got, err := GetAPIKey()
	assert.NoError(t, err)
	assert.Equal(t, "test-key", got)
}

func TestSetupSecrets_MissingAPIKey(t *testing.T) {
	apiKey = "" // Reset global state
	t.Setenv("API_KEY", "") // Explicitly set to empty to override shell environment
	err := SetupSecrets()
	assert.Error(t, err)
}

func TestGetAPIKey_NotInitialized(t *testing.T) {
	apiKey = "" // Reset global state
	_, err := GetAPIKey()
	assert.Error(t, err)
}
