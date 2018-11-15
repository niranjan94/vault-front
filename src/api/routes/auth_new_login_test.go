package routes

import (
	"encoding/json"
	"fmt"
	"github.com/niranjan94/vault-front/src/models"
	testingUtils "github.com/niranjan94/vault-front/src/utils/testing"
	"github.com/niranjan94/vault-front/src/vault"
	assertLib "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewUserLogin(t *testing.T)  {
	assert := assertLib.New(t)
	// Correct username password (Without 2FA)
	_, rec, c := testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: newTestUser.Username,
		Password: newTestUser.Password,
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusOK, rec.Code, "Check if user is able to get 2FA secret for first time")

	otpResponse := models.OtpInformation{}
	json.Unmarshal(rec.Body.Bytes(), &otpResponse)
	assert.NotEmpty(otpResponse.Secret, "Check if user is able to get 2FA secret for first time")
	assert.NotEmpty(otpResponse.Url, "Check if user is able to get 2FA URL for first time")

	// Incorrect username+password (First time)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: "i.do.not.know.who.i.am",
		Password: "wrong_password",
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code, "Check incorrect username+password (w/o 2FA)")

	// Incorrect password (First time)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: newTestUser.Username,
		Password: "wrong_password",
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code, "Check incorrect password (w/o 2FA)")

	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: newTestUser.Username,
		Password: newTestUser.Password,
	})

	// Correct username password (2FA Enabled)
	totpPath := fmt.Sprintf("totp/code/%s", newTestUser.Username)
	otp, err := vault.GetManagerClient().Logical().Read(totpPath)
	if err != nil {
		panic(err)
	}
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: newTestUser.Username,
		Password: newTestUser.Password,
		OTP:      otp.Data["code"].(string),
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusOK, rec.Code, "Check if user is able to login")

	response := LoginResponse{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal("password_expired", response.Status, "Check if user needs to change password")

	waitForNewOtp()

	// Correct username password with new password (2FA Enabled)
	otp, err = vault.GetManagerClient().Logical().Read(totpPath)
	if err != nil {
		panic(err)
	}

	newPassword := "Zxyzzy1234#New"

	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username:    newTestUser.Username,
		Password:    newTestUser.Password,
		NewPassword: newPassword,
		OTP:         otp.Data["code"].(string),
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusCreated, rec.Code, "Check if user password is changed")

	waitForNewOtp()

	otp, err = vault.GetManagerClient().Logical().Read(totpPath)
	if err != nil {
		panic(err)
	}

	// Correct username password (With 2FA)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: newTestUser.Username,
		Password: newPassword,
		OTP:      otp.Data["code"].(string),
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusOK, rec.Code, "Check if user is able to login")

	response = LoginResponse{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NotEmpty(response.Token, "Check if user is able to login and get token")

	// Expire user password
	user := models.NewUser(newTestUser.Username)
	metadata := user.GetMetadata()
	expiredDate :=  metadata.PasswordChangedAt.AddDate(0, -2, 0)
	metadata.PasswordChangedAt = &expiredDate
	user.SetMetadata(metadata)

	waitForNewOtp()

	// Correct username password (2FA Enabled)
	otp, err = vault.GetManagerClient().Logical().Read(totpPath)
	if err != nil {
		panic(err)
	}
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: newTestUser.Username,
		Password: newPassword,
		OTP:      otp.Data["code"].(string),
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusOK, rec.Code, "Check if user is able to login")

	response = LoginResponse{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal("password_expired", response.Status, "Check if user needs to change password")
}