package routes

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/niranjan94/vault-front/src/models"
	"github.com/niranjan94/vault-front/src/utils"
	"github.com/niranjan94/vault-front/src/vault"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

type LoginRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
	OTP         string `json:"otp"`
}

type LoginResponse struct {
	Token  string `json:"token,omitempty"`
	Status string `json:"status,omitempty"`
}

func isPasswordValid(s string) (bool) {
	number := false
	upper := false
	special := false
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
		}
	}
	return number && upper && special && len(s) >= 8
}

var twoFactorPeriod int64 = 30

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		body := new(LoginRequest)
		if err := c.Bind(body); err != nil {
			return err
		}

		body.OTP = strings.TrimSpace(body.OTP)
		body.Username = strings.TrimSpace(body.Username)

		if !isPasswordValid(body.Password) {
			return utils.WriteStatus(c, http.StatusUnauthorized)
		}

		user := models.NewUser(body.Username)
		token := user.Authenticate(body.Password)
		if token == "" {
			return utils.WriteStatus(c, http.StatusUnauthorized)
		}

		if user.HasOTP() {
			isValid, err := user.ValidateOtp(body.OTP)
			if err != nil {
				vault.InvalidateToken(token)
				if strings.Contains(err.Error(), "code already used") {
					return utils.WriteStatus(c, http.StatusConflict)
				}
				log.Println(err.Error())
				return utils.WriteStatus(c, http.StatusInternalServerError)
			}
			if !isValid {
				vault.InvalidateToken(token)
				return utils.WriteStatus(c, http.StatusUnauthorized)
			}
		} else {
			vault.InvalidateToken(token)
			otpInformation, err := user.SetupOTP(twoFactorPeriod)
			if err != nil {
				return utils.WriteStatus(c, http.StatusInternalServerError)
			}
			return c.JSON(http.StatusOK, otpInformation)
		}

		metadata := user.GetMetadata()

		lastValidDateForPassword := time.Now().AddDate(0, -1, 0)

		if body.NewPassword != "" {
			if !isPasswordValid(body.NewPassword) {
				return utils.WriteStatus(c, http.StatusUnprocessableEntity)
			}
			err := user.ChangePassword(token, body.NewPassword)
			vault.InvalidateToken(token)

			if err != nil {
				log.Println(err.Error())
				return utils.WriteStatus(c, http.StatusInternalServerError)
			}

			return utils.WriteStatus(c, http.StatusCreated)
		}

		if metadata.IsNew || (metadata.PasswordChangedAt != nil && metadata.PasswordChangedAt.Before(lastValidDateForPassword)) {
			vault.InvalidateToken(token)
			return c.JSON(http.StatusOK, LoginResponse{
				Status: "password_expired",
			})
		}

		return c.JSON(http.StatusOK, LoginResponse{
			Token: token,
		})
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := vault.InvalidateToken(
			vault.GetTokenFromRequest(c.Request()),
		)
		if err != nil {
			return utils.WriteStatus(c, http.StatusInternalServerError)
		}
		return utils.WriteStatus(c, http.StatusOK)
	}
}
