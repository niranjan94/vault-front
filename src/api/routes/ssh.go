package routes

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/labstack/echo"
	"github.com/niranjan94/vault-front/src/utils"
	"github.com/niranjan94/vault-front/src/vault"
	"log"
	"net/http"
	"path"
	"strings"
)

type sshSigningRequest struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	PublicKey string `json:"publicKey"`
}

type sshSigningResponse struct {
	Username  string      `json:"username"`
	Name      string      `json:"name"`
	Serial    string      `json:"serial"`
	Validity  interface{} `json:"validity"`
	SignedKey string      `json:"signedKey"`
}

func getAllowedRolesFromAllMounts(manager *api.Client, client *api.Client) ([]string, error)  {
	mounts, err := manager.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	var allRoles []string
	var signerPaths []string

	for mountName, mount := range mounts {
		if mount.Type != "ssh" {
			continue
		}

		response, err := manager.Logical().List(path.Join(mountName, "roles"))
		if err != nil {
			return nil, err
		}

		for _, role := range response.Data["keys"].([]interface{}) {
			if val, ok := role.(string); ok {
				signerPaths = append(signerPaths, fmt.Sprintf(path.Join(mountName, "sign", val)))
				allRoles = append(allRoles, val)
			}
		}
	}

	allowedSignerPathsRaw, err := client.Logical().Write("sys/capabilities-self", map[string]interface{}{
		"paths": signerPaths,
	})

	if err != nil {
		return nil, err
	}

	var allowedSignerPaths []string
	for k, v := range allowedSignerPathsRaw.Data {
		if k == "capabilities" {
			continue
		}
		perms := v.([]interface{})
		if perms[0].(string) != "deny" {
			splitPath := strings.Split(k, "/")
			allowedRole := splitPath[len(splitPath)-1]
			allowedMount := splitPath[0]
			if utils.StringSliceContains(allRoles, allowedRole) {
				allowedSignerPaths = append(allowedSignerPaths, allowedMount + ":" + allowedRole)
			}
		}
	}

	return allowedSignerPaths, nil
}

func GetAllowedInstances() echo.HandlerFunc {
	return func(c echo.Context) error {
		manager := vault.GetManagerClient()
		client := vault.GetClientFromContext(c)
		allowedSignerPaths, err := getAllowedRolesFromAllMounts(manager, client)

		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		if allowedSignerPaths == nil {
			allowedSignerPaths = []string{}
		}

		return c.JSON(http.StatusOK, allowedSignerPaths)
	}
}

func GetSignedCertificate() echo.HandlerFunc {
	return func(c echo.Context) error {

		var signingRequest sshSigningRequest
		if err := c.Bind(&signingRequest); err != nil {
			return err
		}

		manager := vault.GetManagerClient()
		client := vault.GetClientFromContext(c)

		requestPayload := map[string]interface{}{
			"public_key": strings.TrimSpace(signingRequest.PublicKey),
		}

		if signingRequest.Username != "" {
			requestPayload["valid_principals"] = signingRequest.Username
		}

		requestRoleRaw := strings.Split(signingRequest.Role, ":")
		if len(requestRoleRaw) < 2 {
			return utils.WriteStatus(c, http.StatusBadRequest)
		}

		credentials, err := client.Logical().Write(
			path.Join(requestRoleRaw[0], "sign", requestRoleRaw[1]),
			requestPayload,
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
			path.Join(requestRoleRaw[0], "roles", requestRoleRaw[1]),
		)

		if err != nil {
			log.Println(err.Error())
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}

		response := sshSigningResponse{
			Username:  requestPayload["valid_principals"].(string),
			Name:      signingRequest.Role,
			Validity:  roleInfo.Data["ttl"],
			SignedKey: credentials.Data["signed_key"].(string),
			Serial:    credentials.Data["serial_number"].(string),
		}

		return c.JSON(http.StatusOK, response)
	}
}
