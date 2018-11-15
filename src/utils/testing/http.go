package testing

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"strings"
)

func NewPostRequestWithBody(body interface{}) (*http.Request, *httptest.ResponseRecorder, echo.Context)  {
	e := echo.New()
	var req *http.Request
	switch t := body.(type) {
	case []byte:
		req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(t))
		break
	case string:
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(t))
		break
	default:
		bodyBytes, err := json.Marshal(t)
		if err != nil {
			panic(err)
		}
		req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return req, rec, c
}

func NewGetRequest() (*http.Request, *httptest.ResponseRecorder, echo.Context)  {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return req, rec, c
}

func NewDeleteRequest() (*http.Request, *httptest.ResponseRecorder, echo.Context)  {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return req, rec, c
}