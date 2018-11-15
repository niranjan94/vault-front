package cmd

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfigForTest(t *testing.T) {
	assert.NoError(t, LoadConfigForTest())
	assert.Equal(t, "/path/to/test.ssl.crt", viper.GetString("api.ssl.crt"))
}

func TestLoadConfig(t *testing.T) {
	assert.Error(t, LoadConfig("invalid-config-path"))
	assert.Error(t, LoadConfig(""))

	assert.NoError(t, LoadConfig(os.Getenv("TEST_WORKING_DIR") + "/config.test.hcl"))
	assert.Equal(t, "/path/to/test.ssl.key", viper.GetString("api.ssl.key"))
}