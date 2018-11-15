package routes

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/niranjan94/vault-front/src/api/auth"
	testingUtils "github.com/niranjan94/vault-front/src/utils/testing"
	"github.com/niranjan94/vault-front/src/vault"
	assertLib "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	assert := assertLib.New(t)

	// Correct username password (2FA Enabled)
	totpPath := fmt.Sprintf("totp/code/%s", twoFactorTestUser.Username)
	otp, err := vault.GetManagerClient().Logical().Read(totpPath)
	if err != nil {
		panic(err)
	}
	_, rec, c := testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: twoFactorTestUser.Username,
		Password: twoFactorTestUser.Password,
		OTP:      otp.Data["code"].(string),
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusOK, rec.Code, "Check if user is able to login")

	response := LoginResponse{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NotEmpty(response.Token, "Check if user is able to login and get token")

	waitForNewOtp()

	// Incorrect password (2FA Enabled)
	// Correct username password (2FA Enabled)
	totpPath = fmt.Sprintf("totp/code/%s", twoFactorTestUser.Username)
	otp, err = vault.GetManagerClient().Logical().Read(totpPath)
	if err != nil {
		panic(err)
	}
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: twoFactorTestUser.Username,
		Password: "wrong_password",
		OTP:      otp.Data["code"].(string),
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code, "Check incorrect password (with 2FA)")

	// Incorrect OTP (2FA Enabled)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: twoFactorTestUser.Username,
		Password: twoFactorTestUser.Password,
		OTP:      "000000",
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code, "Check incorrect OTP")

	// Incorrect OTP (2FA Enabled)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: twoFactorTestUser.Username,
		Password: twoFactorTestUser.Password,
		OTP:      "000000",
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusConflict, rec.Code, "Check incorrect OTP repeat")

	// Missing OTP (2FA Enabled)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: twoFactorTestUser.Username,
		Password: twoFactorTestUser.Password,
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code, "Check missing-required OTP")

	// Incorrect OTP+Password (2FA Enabled)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: twoFactorTestUser.Username,
		Password: "wrong_password",
		OTP:      "000000",
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code, "Check incorrect OTP+Password")

	// Incorrect Username+OTP+Password (2FA Enabled)
	_, rec, c = testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: "i.do.not.know.who.i.am",
		Password: "wrong_password",
		OTP:      "000000",
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code, "Check incorrect Username+OTP+Password")
}

func TestLogout(t *testing.T) {
	assert := assertLib.New(t)

	waitForNewOtp()

	// Correct username password (First time)

	totpPath := fmt.Sprintf("totp/code/%s", twoFactorTestUser.Username)
	otp, err := vault.GetManagerClient().Logical().Read(totpPath)
	if err != nil {
		panic(err)
	}
	_, rec, c := testingUtils.NewPostRequestWithBody(LoginRequest{
		Username: twoFactorTestUser.Username,
		Password: twoFactorTestUser.Password,
		OTP:      otp.Data["code"].(string),
	})

	assert.NoError(Login()(c))
	assert.Equal(http.StatusOK, rec.Code, "Check if user is able to login")

	response := LoginResponse{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NotEmpty(response.Token, "Check if user is able to login and get token")

	req, _, c := testingUtils.NewDeleteRequest()
	h := auth.RequireTokenAuthentication()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Valid key
	authHeader := "Bearer " + response.Token
	req.Header.Set(echo.HeaderAuthorization, authHeader)
	assert.NoError(h(c), "Check if key is valid")

	// Logout user
	assert.NoError(Logout()(c), "Check if user is logged out")

	// Expired key
	he := h(c).(*echo.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code, "Check if logged out user cannot re-use key")
}
