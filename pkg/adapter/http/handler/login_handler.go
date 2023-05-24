package handler

import (
	"fmt"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type Login interface {
	LoginHandler() echo.HandlerFunc
	LogoutHandler() echo.HandlerFunc
}

type LoginHandler struct {
	AuthUseCase usecase.Auth
}

func NewLoginHandler(au usecase.Auth) Login {
	LoginHandler := LoginHandler{
		AuthUseCase: au,
	}
	return &LoginHandler
}

func (l *LoginHandler) LoginHandler() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var fv = &usecase.FormValue{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}

		if err = c.Validate(fv); err != nil {
			return xerrors.Errorf("login validate err: %w", err)
		}

		userId, err := l.AuthUseCase.Login(c, fv)
		if err != nil {
			return fmt.Errorf("login failed err: %w", err)
		}

		return c.JSON(http.StatusOK, echo.Map{
			"userId": userId,
		})
	}
}

func (l *LoginHandler) LogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "already logout",
			})
		}
		l.AuthUseCase.DeleteCookie(c, cookie)
		return c.NoContent(http.StatusOK)
	}
}
