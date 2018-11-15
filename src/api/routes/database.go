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

type CredentialRequest struct {
	Role string `json:"role"`
}

func GetAllowedDatabases() echo.HandlerFunc {
	return func(c echo.Context) error {
		manager := vault.GetManagerClient()
		client := vault.GetClientFromContext(c)

		response, err := manager.Logical().List("database/roles")
		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		var databaseRolePaths []string

		for _, role := range response.Data["keys"].([]interface{}) {
			databaseRolePaths = append(databaseRolePaths, fmt.Sprintf("database/creds/%s", role))
		}

		allowedRolesRaw, err := client.Logical().Write("sys/capabilities-self", map[string]interface{}{
			"paths": databaseRolePaths,
		})

		var allowedRoles []string

		for k, v := range allowedRolesRaw.Data {
			perms := v.([]interface{})
			if perms[0].(string) != "deny" {
				splitPath := strings.Split(k, "/")
				allowedRoles = append(allowedRoles, splitPath[len(splitPath) - 1])
			}
		}

		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, allowedRoles)
	}
}

func GetCredential() echo.HandlerFunc {
	return func(c echo.Context) error {

		var credentialRequest CredentialRequest
		if err := c.Bind(&credentialRequest); err != nil {
			return err
		}

		client := vault.GetClientFromContext(c)
		credentials, err := client.Logical().Read(
			fmt.Sprintf("database/creds/%s", credentialRequest.Role),
		)

		if err != nil {
			errorMessage := err.Error()
			if strings.Contains(errorMessage, "permission denied") {
				return utils.WriteStatus(c, http.StatusUnauthorized)
			}
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}


		return c.JSON(http.StatusOK, credentials.Data)
	}
}
