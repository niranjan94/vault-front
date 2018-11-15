package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func GetManagerClient() (*api.Client) {
	if token := viper.GetString("vault.token"); token != "" {
		return GetClient(token)
	}
	return GetBaseClient()
}