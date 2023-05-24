package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Csrf interface {
	CsrfHandler() echo.HandlerFunc
}

type CsrfHandler struct {
}

func NewCsrfHandler() Csrf {
	CsrfHandler := CsrfHandler{}
	return &CsrfHandler
}

func (l *CsrfHandler) CsrfHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "csrf-token set")
	}
}
