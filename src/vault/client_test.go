package vault

import (
	"github.com/labstack/echo"
	testingUtils "github.com/niranjan94/vault-front/src/utils/testing"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetClient(t *testing.T) {
	token := viper.GetString("vault.token")
	client := GetClient(token)
	assert.Equal(t, token, client.Token())
}

func TestGetClientFromContext(t *testing.T) {
	token := viper.GetString("vault.token")
	req, _, c := testingUtils.NewGetRequest()
	auth := "Bearer " + token
	req.Header.Set(echo.HeaderAuthorization, auth)

	client := GetClientFromContext(c)

	assert.Equal(t, token, client.Token())
}

func TestValidateToken(t *testing.T) {
	assert.True(t,  ValidateToken(viper.GetString("vault.token")))
	assert.False(t,  ValidateToken("not-a-real-token"))
}

func TestInvalidateToken(t *testing.T)  {
	token := os.Getenv("PII_TEST_VAULT_TOKEN")

	assert.True(t,  ValidateToken(token))
	assert.NoError(t, InvalidateToken(token))
	assert.False(t,  ValidateToken(token))
}