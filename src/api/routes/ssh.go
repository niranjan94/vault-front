package routes

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/niranjan94/vault-front/src/utils"
	"github.com/niranjan94/vault-front/src/vault"
	"log"
	"net/http"
	"strings"
)

type SigningRequest struct {
	Role      string `json:"role"`
	PublicKey string `json:"publicKey"`
}

type SigningResponse struct {
	Username  string      `json:"username"`
	Name      string      `json:"name"`
	Validity  interface{} `json:"validity"`
	SignedKey string      `json:"signedKey"`
}

func GetAllowedInstances() echo.HandlerFunc {
	return func(c echo.Context) error {
		manager := vault.GetManagerClient()
		client := vault.GetClientFromContext(c)

		response, err := manager.Logical().List("ssh-client-signer/roles")
		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		var allRoles []string
		var signerPaths []string

		for _, role := range response.Data["keys"].([]interface{}) {
			signerPaths = append(signerPaths, fmt.Sprintf("ssh-client-signer/sign/%s", role))
			allRoles = append(allRoles, fmt.Sprintf("%s", role))
		}

		allowedSignerPathsRaw, err := client.Logical().Write("sys/capabilities-self", map[string]interface{}{
			"paths": signerPaths,
		})
		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		var allowedSignerPaths []string

		for k, v := range allowedSignerPathsRaw.Data {
			perms := v.([]interface{})
			if perms[0].(string) != "deny" {
				splitPath := strings.Split(k, "/")
				allowedRole := splitPath[len(splitPath)-1]
				if utils.StringSliceContains(allRoles, allowedRole) {
					allowedSignerPaths = append(allowedSignerPaths, allowedRole)
				}
			}
		}

		return c.JSON(http.StatusOK, allowedSignerPaths)
	}
}

func GetSignedCertificate() echo.HandlerFunc {
	return func(c echo.Context) error {

		var signingRequest SigningRequest
		if err := c.Bind(&signingRequest); err != nil {
			return err
		}

		manager := vault.GetManagerClient()
		client := vault.GetClientFromContext(c)
		credentials, err := client.Logical().Write(
			fmt.Sprintf("ssh-client-signer/sign/%s", signingRequest.Role),
			map[string]interface{}{
				"public_key": strings.TrimSpace(signingRequest.PublicKey),
			},
		)

		if err != nil {
			errorMessage := err.Error()
			if strings.Contains(errorMessage, "permission denied") {
				return utils.WriteStatus(c, http.StatusUnauthorized)
			}
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		roleInfo, err := manager.Logical().Read(
			fmt.Sprintf("ssh-client-signer/roles/%s", signingRequest.Role),
		)

		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		response := SigningResponse{
			Username:  roleInfo.Data["default_user"].(string),
			Name:      signingRequest.Role,
			Validity:  roleInfo.Data["ttl"],
			SignedKey: credentials.Data["signed_key"].(string),
		}

		return c.JSON(http.StatusOK, response)
	}
}
