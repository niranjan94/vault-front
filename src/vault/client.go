package vault

import (
	"bytes"
	"context"
	"github.com/hashicorp/vault/api"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func GetBaseClient() (*api.Client)  {
	vaultClient, err := api.NewClient(&api.Config{
		Address: viper.GetString("vault.address"),
	})
	if err != nil {
		panic(err)
	}
	return vaultClient
}

func GetClient(token string) (*api.Client) {
	vaultClient := GetBaseClient()
	if token != "" {
		vaultClient.SetToken(token)
	}
	return vaultClient
}

func GetClientFromContext(c echo.Context) (*api.Client) {
	return GetClientFromRequest(c.Request())
}

func GetClientFromRequest(r *http.Request) (*api.Client) {
	return GetClient(GetTokenFromRequest(r))
}

func GetTokenFromRequest(r *http.Request) string  {
	token := r.Header.Get("Authorization")
	return token[len("Bearer") + 1:]
}

func ValidateToken(token string) (bool)  {
	vaultClient := GetClient(token)
	request := vaultClient.NewRequest("POST", "/v1/sys/tools/random")
	response, err := vaultClient.RawRequestWithContext(context.Background(), request)
	if err != nil {
		return false
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return false
	}
	return true
}

func InvalidateToken(token string) (error)  {
	client := GetClient(token)
	request := client.NewRequest("POST", "/v1/auth/token/revoke-self")
	_, err := client.RawRequestWithContext(context.Background(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}