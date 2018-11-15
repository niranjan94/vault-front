package routes

import (
	"fmt"
	"github.com/niranjan94/vault-front/src/cmd"
	"github.com/niranjan94/vault-front/src/models"
	"github.com/niranjan94/vault-front/src/vault"
	"time"
)

var (
	token string
	restrictedToken string
	twoFactorTestUser *LoginRequest
	newTestUser *LoginRequest
	plainTestUser *LoginRequest
)

func init()  {

	twoFactorPeriod = 1
	cmd.LoadConfigForTest()
	token = createAndLoginTestUser(true)
	restrictedToken = createAndLoginTestUser(false)
	cmd.LoadConfigForTest()
	twoFactorTestUser = &LoginRequest{
		Username: "john.doe",
		Password: "Zxyzzy1234#",
	}
	plainTestUser = &LoginRequest{
		Username: "jane.doe",
		Password: "Zxyzzy1234#",
	}
	newTestUser = &LoginRequest{
		Username: "new.jane.doe",
		Password: "Zxyzzy1234#",
	}
	createTestUser(twoFactorTestUser, true, false, true)
	createTestUser(plainTestUser, false, false, true)
	createTestUser(newTestUser, false, true, false)
}

func waitForNewOtp()  {
	time.Sleep(time.Duration(twoFactorPeriod) * time.Second)
}

func createAndLoginTestUser(allDatabases bool) string {
	primaryTestUser := &LoginRequest{
		Username: fmt.Sprintf("%d-john.jane.doe", time.Now().UnixNano()),
		Password: "Zxyzzy1234#",
	}
	createTestUser(primaryTestUser, true, false, allDatabases)
	client := vault.GetManagerClient()
	userpassPath := fmt.Sprintf("auth/userpass/login/%s", primaryTestUser.Username)
	data, err := client.Logical().Write(userpassPath, map[string]interface{}{
		"password": primaryTestUser.Password,
	})
	if err != nil {
		panic(err)
	}

	return data.Auth.ClientToken
}

func createTestUser(body *LoginRequest, withTwoFactor bool, isNew bool, allDatabases bool) {
	client := vault.GetManagerClient()

	userpassPath := fmt.Sprintf("auth/userpass/users/%s", body.Username)
	totpPath := fmt.Sprintf("totp/keys/%s", body.Username)

	if withTwoFactor {
		_, err := client.Logical().Write(totpPath, map[string]interface{}{
			"generate":     true,
			"issuer":       "Vault Front",
			"account_name": body.Username,
			"period"    : twoFactorPeriod,
		})
		if err != nil {
			panic(err)
		}
	}

	policies := "default,test-user-policy-restricted"

	if allDatabases {
		policies = "default,test-user-policy"
	}

	_, err := client.Logical().Write(userpassPath, map[string]interface{}{
		"password": body.Password,
		"policies": policies,
	})

	if err != nil {
		panic(err)
	}

	now := time.Now()

	user := models.NewUser(body.Username)
	metadata := user.GetMetadata()
	metadata.IsNew = isNew
	metadata.LastLoginAt = &now
	metadata.PasswordChangedAt = &now
	err = user.SetMetadata(metadata)
	if err != nil {
		panic(err)
	}
}

