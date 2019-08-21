package routes

import (
	"github.com/hashicorp/vault/api"
	"github.com/labstack/echo"
	"github.com/niranjan94/vault-front/src/utils"
	"github.com/niranjan94/vault-front/src/vault"
	"github.com/spf13/cast"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type systemCredentialRequest struct {
	Role string `json:"role"`
}

type systemCredentialResponse struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Validity      int    `json:"validity"`
}

func getSystemCredentialPaths(manager *api.Client) ([]string, error) {
	var roles []string
	metadataPath := "instance-creds/metadata"
	metadata, err := manager.Logical().List(metadataPath)
	if err != nil {
		return nil, err
	}
	for _, subKey := range cast.ToStringSlice(metadata.Data["keys"]) {
		instances, err := manager.Logical().List(path.Join(metadataPath, subKey))
		if err != nil {
			return nil, err
		}
		for _, instance := range cast.ToStringSlice(instances.Data["keys"]) {
			users, err := manager.Logical().List(path.Join(metadataPath, subKey, instance))
			if err != nil {
				return nil, err
			}
			for _, user := range cast.ToStringSlice(users.Data["keys"]) {
				roles = append(roles, path.Join("instance-creds", "data", subKey, instance, user))
			}
		}
	}
	return roles, nil
}

func GetAllowedSystemCredentials() echo.HandlerFunc {
	return func(c echo.Context) error {
		manager := vault.GetManagerClient()
		client := vault.GetClientFromContext(c)

		rolePaths, err := getSystemCredentialPaths(manager)
		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		allowedRolesRaw, err := client.Logical().Write("sys/capabilities-self", map[string]interface{}{
			"paths": rolePaths,
		})

		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		var allowedRoles []string

		for k, v := range allowedRolesRaw.Data {
			perms := v.([]interface{})
			if perms[0].(string) != "deny" && utils.StringSliceContains(rolePaths, k) {
				allowedRoles = append(allowedRoles, k)
			}
		}

		if allowedRoles == nil {
			allowedRoles = []string{}
		}

		return c.JSON(http.StatusOK, allowedRoles)
	}
}

func GetSystemCredential() echo.HandlerFunc {
	return func(c echo.Context) error {

		var credentialRequest systemCredentialRequest
		if err := c.Bind(&credentialRequest); err != nil {
			return err
		}

		if !strings.HasPrefix(credentialRequest.Role, "instance-creds/data") {
			return utils.WriteStatus(c, http.StatusBadRequest)
		}

		client := vault.GetClientFromContext(c)
		credentials, err := client.Logical().Read(credentialRequest.Role)

		if err != nil {
			errorMessage := err.Error()
			if strings.Contains(errorMessage, "permission denied") {
				return utils.WriteStatus(c, http.StatusUnauthorized)
			}
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		data := credentials.Data["data"].(map[string]interface{})
		metadata := credentials.Data["metadata"].(map[string]interface{})

		createdAt, err := time.Parse(
			time.RFC3339,
			metadata["created_time"].(string),
		)
		if err != nil {
			return err
		}
		durationRemaining := createdAt.Add(time.Hour * 4).Sub(time.Now())
		response := systemCredentialResponse{
			Username: data["username"].(string),
			Password: data["password"].(string),
			Validity: int(durationRemaining.Seconds()),
		}
		return c.JSON(http.StatusOK, response)
	}
}
