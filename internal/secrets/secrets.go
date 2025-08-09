package secrets

import (
	"errors"
	"os"
)

var apiKey string

// SetupSecrets loads the API key from the environment and caches it.
// This function should be called once during application startup.
func SetupSecrets() error {
	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		return errors.New("API_KEY environment variable is not set")
	}
	return nil
}

// GetAPIKey retrieves the cached API key.
func GetAPIKey() (string, error) {
	if apiKey == "" {
		return "", errors.New("API key not initialized")
	}
	return apiKey, nil
}
