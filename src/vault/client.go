package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func GetBaseClient() *api.Client {
	vaultClient, err := api.NewClient(&api.Config{
		Address: viper.GetString("vault.address"),
	})
	if err != nil {
		panic(err)
	}
	return vaultClient
}

func GetClient(token string) *api.Client {
	vaultClient := GetBaseClient()
	if token != "" {
		vaultClient.SetToken(token)
	}
	return vaultClient
}

func GetClientFromContext(c echo.Context) *api.Client {
	return GetClientFromRequest(c.Request())
}

func GetClientFromRequest(r *http.Request) *api.Client {
	return GetClient(GetTokenFromRequest(r))
}

func GetTokenFromRequest(r *http.Request) string  {
	token := r.Header.Get("Authorization")
	return token[len("Bearer") + 1:]
}

func ValidateToken(token string) bool {
	vaultClient := GetClient(token)
	self, err := vaultClient.Auth().Token().LookupSelf()
	if err != nil {
		return false
	}
	if self.Data == nil {
		return false
	}
	return true
}

func InvalidateToken(token string) error {
	client := GetClient(token)
	err := client.Auth().Token().RevokeSelf(token)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}