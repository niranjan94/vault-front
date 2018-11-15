package auth

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/niranjan94/vault-front/src/vault"
)

func TokenValidator() middleware.KeyAuthValidator {
	return func(token string, c echo.Context) (bool, error) {
		if vault.ValidateToken(token) {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func RequireTokenAuthentication() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
		Validator: TokenValidator(),
	})
}