package routes

import (
	"fmt"
	"github.com/niranjan94/vault-front/src/cmd"
	"github.com/niranjan94/vault-front/src/models"
	"github.com/niranjan94/vault-front/src/utils"
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
		Email: "john.doe@vault.com",
		Password: "Zxyzzy1234#",
	}
	plainTestUser = &LoginRequest{
		Email: "jane.doe@vault.com",
		Password: "Zxyzzy1234#",
	}
	newTestUser = &LoginRequest{
		Email: "new.jane.doe@vault.com",
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
		Email: fmt.Sprintf("%d-john.jane.doe@vault.com", time.Now().UnixNano()),
		Password: "Zxyzzy1234#",
	}
	createTestUser(primaryTestUser, true, false, allDatabases)
	client := vault.GetManagerClient()
	id := utils.SHA512(primaryTestUser.Email)
	userpassPath := fmt.Sprintf("auth/userpass/login/%s", id)
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
	email := utils.SHA512(body.Email)

	userpassPath := fmt.Sprintf("auth/userpass/users/%s", email)
	totpPath := fmt.Sprintf("totp/keys/%s", email)

	if withTwoFactor {
		_, err := client.Logical().Write(totpPath, map[string]interface{}{
			"generate":     true,
			"issuer":       "Vault Front",
			"account_name": body.Email,
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

	user := models.NewUser(body.Email)
	metadata := user.GetMetadata()
	metadata.IsNew = isNew
	metadata.LastLoginAt = &now
	metadata.PasswordChangedAt = &now
	err = user.SetMetadata(metadata)
	if err != nil {
		panic(err)
	}
}

