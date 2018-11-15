package vault

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetManagerClient(t *testing.T) {
	assert.Equal(t, viper.GetString("vault.token"), GetManagerClient().Token())

	tokenBackup := viper.GetString("vault.token")
	viper.Set("vault.token", "")

	assert.Empty(t, GetManagerClient().Token())

	viper.Set("vault.token", tokenBackup)
}