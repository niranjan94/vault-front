package models

import (
	"fmt"
	"github.com/niranjan94/vault-front/src/vault"
	"log"
	"net/url"
	"time"
)

type User struct {
	Username string

	userpassPath string
	totpKeyPath  string
	totpCodePath string
	metadataPath string
}

type UserMetadata struct {
	IsNew             bool
	LastLoginAt       *time.Time
	PasswordChangedAt *time.Time
}

type OtpInformation struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
}

func (u *User) Authenticate(password string) (string) {
	client := vault.GetManagerClient()
	data, err := client.Logical().Write(u.userpassPath, map[string]interface{}{
		"password": password,
	})
	if err != nil || data == nil || data.Auth == nil {
		return ""
	}
	return data.Auth.ClientToken
}

func (u *User) ChangePassword(token string, password string) error {
	userpassPasswordPath := fmt.Sprintf("auth/userpass/users/%s/password", u.Username)
	client := vault.GetManagerClient()
	_, err := client.Logical().Write(userpassPasswordPath, map[string]interface{}{
		"password": password,
	})

	if err != nil {
		return err
	}

	currentTime := time.Now()

	metadata := u.GetMetadata()
	metadata.IsNew = false
	metadata.PasswordChangedAt = &currentTime
	metadata.LastLoginAt = &currentTime

	return u.SetMetadata(metadata)
}

func (u *User) HasOTP() bool {
	client := vault.GetManagerClient()
	otpKey, err := client.Logical().Read(u.totpKeyPath)
	return err == nil && otpKey != nil
}

func (u *User) ValidateOtp(otp string) (bool, error) {
	if otp == "" {
		return false, nil
	}

	client := vault.GetManagerClient()

	validation, err := client.Logical().Write(
		u.totpCodePath,
		map[string]interface{}{
			"code": otp,
		},
	)
	if err != nil {
		return false, err
	}
	return validation.Data["valid"].(bool), nil

}

func (u *User) SetupOTP(twoFactorPeriod int64) (*OtpInformation, error) {
	client := vault.GetManagerClient()
	client.Logical().Delete(u.totpKeyPath)
	otpInfo, err := client.Logical().Write(u.totpKeyPath, map[string]interface{}{
		"generate":     true,
		"issuer":       "Vault Front",
		"account_name": u.Username,
		"period": 		twoFactorPeriod,
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	totpUrl, err := url.Parse(otpInfo.Data["url"].(string))

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &OtpInformation{
		Url:    totpUrl.String(),
		Secret: totpUrl.Query().Get("secret"),
	}, nil

}

func (u *User) GetMetadata() *UserMetadata {
	client := vault.GetManagerClient()
	metadata, err := client.Logical().Read(u.metadataPath)
	if err == nil && metadata != nil && metadata.Data != nil{
		lastLoginAt, _ := time.Parse(time.RFC3339, metadata.Data["last_login_at"].(string))
		passwordChangedAt, _ := time.Parse(time.RFC3339, metadata.Data["password_changed_at"].(string))
		return &UserMetadata{
			IsNew:             metadata.Data["is_new"].(bool),
			LastLoginAt:       &lastLoginAt,
			PasswordChangedAt: &passwordChangedAt,
		}
	} else {
		return &UserMetadata{
			IsNew:             true,
			LastLoginAt:       nil,
			PasswordChangedAt: nil,
		}
	}
}

func (u *User) SetMetadata(metadata *UserMetadata) error {
	client := vault.GetManagerClient()
	_, err := client.Logical().Write(u.metadataPath, map[string]interface{}{
		"is_new":              metadata.IsNew,
		"last_login_at":       metadata.LastLoginAt.Format(time.RFC3339),
		"password_changed_at": metadata.PasswordChangedAt.Format(time.RFC3339),
	})
	return err
}

func NewUser(username string) *User {
	user := User{
		Username:     username,
		metadataPath: fmt.Sprintf("kv/metadata/auth/userpass/%s", username),
		userpassPath: fmt.Sprintf("auth/userpass/login/%s", username),
		totpKeyPath:  fmt.Sprintf("totp/keys/%s", username),
		totpCodePath: fmt.Sprintf("totp/code/%s", username),
	}

	return &user
}
