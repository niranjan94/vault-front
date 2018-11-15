package auth

import (
	"github.com/labstack/echo"
	"github.com/niranjan94/vault-front/src/cmd"
	testingUtils "github.com/niranjan94/vault-front/src/utils/testing"
	"github.com/spf13/viper"
	assertLib "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)


func init() {
	cmd.LoadConfigForTest()
}

func TestRequireTokenAuthentication(t *testing.T) {
	req, _, c := testingUtils.NewGetRequest()
	assert := assertLib.New(t)

	h := RequireTokenAuthentication()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Valid key
	auth := "Bearer " + viper.GetString("vault.token")
	req.Header.Set(echo.HeaderAuthorization, auth)
	assert.NoError(h(c))

	// Invalid key
	auth = "Bearer invalid-key"
	req.Header.Set(echo.HeaderAuthorization, auth)
	he := h(c).(*echo.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)

	// Missing Authorization header
	req.Header.Del(echo.HeaderAuthorization)
	he = h(c).(*echo.HTTPError)
	assert.Equal(http.StatusBadRequest, he.Code)
}
