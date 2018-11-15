package main

import (
	"flag"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/hashicorp/vault/helper/mlock"
	"github.com/niranjan94/vault-front/src/api"
	"github.com/niranjan94/vault-front/src/cmd"
	"github.com/spf13/viper"
	"log"
	"os"
)

var config string

func init() {
	flag.StringVar(&config, "config", "", "Path to config.hcl")
	flag.Parse()
}

const mlockError = `
Failed to use mlock to prevent swap usage: %s

vault-front uses mlock similar to Vault. See here for details:
https://www.vaultproject.io/docs/configuration/index.html#disable_mlock

To enable mlock without launching vault-front as root:
sudo setcap cap_ipc_lock=+ep $(readlink -f $(which vault-front))

To disable mlock entirely, set disableMLock to "true" in config file
`

func main() {
	err := cmd.LoadConfig(config)

	if err != nil {
		fmt.Printf("Failed to load config file. => %s\n", err.Error())
		os.Exit(-1)
	}

	if !viper.GetBool("disableMLock") {
		if err := mlock.LockMemory(); err != nil {
			log.Fatalf(mlockError, err.Error())
		}
	}

	api.StartApiServer(
		rice.MustFindBox("ui/dist").HTTPBox(),
		true,
	)
}
