package utils

import (
	"github.com/labstack/echo"
	assertLib "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteStatusSuccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert := assertLib.New(t)
	WriteStatus(c, http.StatusOK)
	assert.Equal(http.StatusOK, rec.Code)
}

func TestWriteStatusUnauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert := assertLib.New(t)
	WriteStatus(c, http.StatusUnauthorized)
	assert.Equal(http.StatusUnauthorized, rec.Code)
}


func TestWriteStatusError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert := assertLib.New(t)
	WriteStatus(c, http.StatusInternalServerError)
	assert.Equal(http.StatusInternalServerError, rec.Code)
}