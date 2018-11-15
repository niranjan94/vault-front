package routes

import (
	"encoding/json"
	"github.com/labstack/echo"
	testingUtils "github.com/niranjan94/vault-front/src/utils/testing"
	assertLib "github.com/stretchr/testify/assert"
	"net/http"
	"sort"
	"testing"
)

type credentialResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestGetAllowedDatabases(t *testing.T) {
	assert := assertLib.New(t)
	req, rec, c := testingUtils.NewGetRequest()
	req.Header.Set(echo.HeaderAuthorization, "Bearer " + token)
	assert.NoError(GetAllowedDatabases()(c))
	assert.Equal(http.StatusOK, rec.Code)

	allRoles := []string{"database-role-one", "database-role-two", "database-role-three"}
	sort.Strings(allRoles)

	restrictedRoles := []string{"database-role-one"}

	var allowedRoles []string
	assert.NoError(json.Unmarshal(rec.Body.Bytes(), &allowedRoles))
	sort.Strings(allowedRoles)
	assert.Equal(allRoles, allowedRoles)

	req, rec, c = testingUtils.NewGetRequest()
	req.Header.Set(echo.HeaderAuthorization, "Bearer " + restrictedToken)
	assert.NoError(GetAllowedDatabases()(c))
	assert.Equal(http.StatusOK, rec.Code)

	assert.NoError(json.Unmarshal(rec.Body.Bytes(), &allowedRoles))
	sort.Strings(allowedRoles)
	assert.Equal(restrictedRoles, allowedRoles)
}

func TestGetCredential(t *testing.T)  {
	assert := assertLib.New(t)
	req, rec, c := testingUtils.NewPostRequestWithBody(CredentialRequest{
		Role: "database-role-one",
	})
	req.Header.Set(echo.HeaderAuthorization, "Bearer " + token)
	assert.NoError(GetCredential()(c))
	assert.Equal(http.StatusOK, rec.Code)

	var credentialResponse credentialResponse
	assert.NoError(json.Unmarshal(rec.Body.Bytes(), &credentialResponse))
	assert.Equal(credentialResponse.Username, "john-doe")
	assert.Equal(credentialResponse.Password, "xyzzy1")

	req, rec, c = testingUtils.NewPostRequestWithBody(CredentialRequest{
		Role: "database-role-two",
	})
	req.Header.Set(echo.HeaderAuthorization, "Bearer " + restrictedToken)
	assert.NoError(GetCredential()(c))
	assert.Equal(http.StatusUnauthorized, rec.Code)
}