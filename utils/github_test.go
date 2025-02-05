package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the AppClient function
func TestAppClient(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("GITHUB_PRIVATE_KEY", "test_private_key")
	os.Setenv("GITHUB_APP_ID", "123456")
	os.Setenv("GITHUB_INSTALLATION_ID", "654321")

	client, err := AppClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestAppClientInvalidAppID(t *testing.T) {
	// Set up environment variables with invalid app ID
	os.Setenv("GITHUB_PRIVATE_KEY", "test_private_key")
	os.Setenv("GITHUB_APP_ID", "invalid_app_id")
	os.Setenv("GITHUB_INSTALLATION_ID", "654321")

	client, err := AppClient()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestAppClientInvalidInstallationID(t *testing.T) {
	// Set up environment variables with invalid installation ID
	os.Setenv("GITHUB_PRIVATE_KEY", "test_private_key")
	os.Setenv("GITHUB_APP_ID", "123456")
	os.Setenv("GITHUB_INSTALLATION_ID", "invalid_installation_id")

	client, err := AppClient()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestAppClientInvalidPrivateKey(t *testing.T) {
	// Set up environment variables with invalid private key
	os.Setenv("GITHUB_PRIVATE_KEY", "")
	os.Setenv("GITHUB_APP_ID", "123456")
	os.Setenv("GITHUB_INSTALLATION_ID", "654321")

	client, err := AppClient()
	assert.Error(t, err)
	assert.Nil(t, client)
}
