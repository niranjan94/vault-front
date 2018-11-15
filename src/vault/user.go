package vault

import (
	"github.com/hashicorp/vault/api"
)

func GetSelf(client *api.Client) map[string]interface{} {
	self, err := client.Logical().Read("auth/token/lookup-self")
	if err != nil {
		panic(err)
	}
	return self.Data
}

func GetSelfEntityId(client *api.Client) string {
	self := GetSelf(client)
	return self["entity_id"].(string)
}