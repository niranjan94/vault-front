package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/matryer/resync"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var once resync.Once
var managerClient *api.Client

func getInstanceIdentityInfo(param string) (string, error) {
	metadataEndpoint := viper.GetString("vault.instanceMetadataService")
	if strings.HasSuffix(metadataEndpoint, "/") {
		metadataEndpoint = strings.TrimSuffix(metadataEndpoint, "/")
	}
	resp, err := http.Get(metadataEndpoint + "/latest/dynamic/instance-identity/" + param)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func authenticateViaEC2() (*api.Secret, error) {
	pkcs7, err := getInstanceIdentityInfo("pkcs7")
	if err != nil {
		return nil, err
	}
	pkcs7 = strings.Replace(pkcs7, "\n", "", -1)
	return GetBaseClient().Logical().Write("auth/aws/login", map[string]interface{}{
		"role":  viper.GetString("vault.authRole"),
		"pkcs7": pkcs7,
		"nonce": viper.GetString("vault.authNonce"),
	})
}

func GetManagerClient() (*api.Client) {
	once.Do(func() {
		token := viper.GetString("vault.token");
		switch viper.GetString("vault.authMode") {
		case "ec2":
			{
				secret, err := authenticateViaEC2()
				if err != nil {
					log.Println(err.Error())
					once.Reset()
					break
				}
				token = secret.Auth.ClientToken
				time.AfterFunc(time.Duration(secret.Auth.LeaseDuration - 10) * time.Second, func() {
					once.Reset()
				})
				break
			}
		}
		managerClient = GetClient(token)
	})

	return managerClient
}
