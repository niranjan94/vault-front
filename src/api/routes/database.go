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

type databaseCredentialRequest struct {
	Role string `json:"role"`
}

type databaseCredentialResponse struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ConnectionUrl string `json:"connectionUrl"`
	Validity      int    `json:"validity"`
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
		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		var allowedRoles []string

		for k, v := range allowedRolesRaw.Data {
			perms := v.([]interface{})
			if perms[0].(string) != "deny" {
				splitPath := strings.Split(k, "/")
				allowedRoles = append(allowedRoles, splitPath[len(splitPath)-1])
			}
		}

		if allowedRoles == nil {
			allowedRoles = []string{}
		}

		return c.JSON(http.StatusOK, allowedRoles)
	}
}

func GetDatabaseCredential() echo.HandlerFunc {
	return func(c echo.Context) error {

		var credentialRequest databaseCredentialRequest
		if err := c.Bind(&credentialRequest); err != nil {
			return err
		}

		manager := vault.GetManagerClient()
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

		roleInfo, err := manager.Logical().Read(
			fmt.Sprintf("database/roles/%s", credentialRequest.Role),
		)

		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		databaseInfo, err := manager.Logical().Read(
			fmt.Sprintf("database/config/%s", roleInfo.Data["db_name"]),
		)

		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		connectionUrl := databaseInfo.Data["connection_details"].(map[string]interface{})["connection_url"].(string)

		response := databaseCredentialResponse{
			Username: credentials.Data["username"].(string),
			Password: credentials.Data["password"].(string),
			Validity: credentials.LeaseDuration,
		}

		connectionUrl = strings.Replace(connectionUrl, "{{username}}", response.Username, 1)
		connectionUrl = strings.Replace(connectionUrl, "{{password}}", response.Password, 1)

		response.ConnectionUrl = connectionUrl

		return c.JSON(http.StatusOK, response)
	}
}
