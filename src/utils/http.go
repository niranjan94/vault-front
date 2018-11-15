package utils

import (
	"github.com/labstack/echo"
	"net/http"
)

type StatusResponse struct {
	Message string `json:"status,omitempty"`
}

func WriteStatus(c echo.Context, code int) error {
	return c.JSON(code, StatusResponse{
		Message: http.StatusText(code),
	})
}
