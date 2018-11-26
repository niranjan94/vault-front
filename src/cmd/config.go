package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/hcl"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

type config struct {
	Api struct {
		ListenAddress string `json:"listenAddress"`
		SSL           struct {
			Enabled bool   `json:"enabled"`
			Crt     string `json:"crt"`
			Key     string `json:"key"`
		} `json:"ssl"`
	} `json:"api"`

	CORS struct {
		AllowedOrigins []string `json:"allowedOrigins"`
	} `json:"cors"`

	DisableMLock bool `json:"disableMLock"`

	Vault struct {
		Address  string `json:"address"`
		AuthMode string `json:"authMode"`
		AuthRole string `json:"authRole"`
		Token    string `json:"token"`
	} `json:"vault"`
}

func loadFromHCL(pathToFile string) error {
	var c config
	fileContent, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		return err
	}

	hcl.Unmarshal(fileContent, &c)
	jsonData, _ := json.Marshal(c)
	viper.SetConfigType("json")
	err = viper.ReadConfig(bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	return nil
}

func LoadConfigForTest() error {
	return LoadConfig(os.Getenv("TEST_WORKING_DIR") + "/config.test.hcl")
}

func LoadConfig(pathToConfig string) error {
	replacer := strings.NewReplacer(".", "_")

	var err error

	if pathToConfig != "" {
		err = loadFromHCL(pathToConfig)
	} else {
		err = loadFromHCL("config.hcl")
	}

	if err != nil {
		return err
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PII")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("api.listenAddress", "0.0.0.0:8000")
	viper.SetDefault("api.ssl.enabled", false)
	viper.SetDefault("vault.address", "http://127.0.0.1:8200")
	viper.SetDefault("vault.instanceMetadataService", "http://169.254.169.254")
	viper.SetDefault("disableMLock", false)

	return nil
}
