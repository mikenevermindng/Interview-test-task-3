package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("APP_ENV", "test")
	os.Chdir("../")

	// Call the Initialize function
	Initialize()

	// Call the NewConfiguration function
	config := NewConfiguration()

	// Check if the configuration values are properly initialized
	assert.NotNil(t, config.Monitor)
	assert.Equal(t, 100, config.Monitor.MaxMonitorConcurrency)
	assert.Equal(t, 6000, config.Monitor.MonitorIntervalMs)
	assert.Equal(t, 5, config.Monitor.MaxRedirect)
	assert.Equal(t, 5000, config.Monitor.RequestTimeout)
	assert.Equal(t, "wz$Z@%%oixtJ&^w&8Bae6^", config.Api.AdminSecret)
	assert.Equal(t, 8080, config.Api.Port)
	assert.Equal(t, "test.db", config.Database.Sqlite)

}
