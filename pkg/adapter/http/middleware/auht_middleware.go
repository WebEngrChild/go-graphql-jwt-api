package middleware

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type AuthMiddleware struct {
	AuthUseCase usecase.Auth
}

func NewAuthMiddleware(ju usecase.Auth) Auth {
	AuthMiddleware := AuthMiddleware{
		AuthUseCase: ju,
	}
	return &AuthMiddleware
}

func (j *AuthMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if j.isSkippedPath(c.Request().URL.Path, c.Request().Referer()) {
			if err := next(c); err != nil {
				return xerrors.Errorf("AuthMiddleware error path: %s: %w", c.Request().URL.Path, err)
			}
			return nil
		}

		cookie, err := c.Cookie("token")
		if err != nil {
			return xerrors.Errorf("AuthMiddleware not extract cookie: %w", err)
		}

		claims, err := j.AuthUseCase.JwtParser(cookie.Value)
		if err != nil {
			j.AuthUseCase.DeleteCookie(c, cookie)
			return fmt.Errorf("failed to parse jwt claims: %w", err)
		}

		var cl = *claims
		uId := cl["user_id"].(string)
		if err := j.AuthUseCase.IdentifyJwtUser(uId); err != nil {
			j.AuthUseCase.DeleteCookie(c, cookie)
			return fmt.Errorf("failed to personal authentication: %w", err)
		}

		if err := next(c); err != nil {
			return xerrors.Errorf("failed to AuthMiddleware err: %w", err)
		}

		return nil
	}
}

func (j *AuthMiddleware) isSkippedPath(reqPath, refPath string) bool {
	skippedPaths := []string{"/healthcheck", "/csrf-cookie", "/login", "/logout", "/playground"}
	for _, path := range skippedPaths {
		if strings.Contains(reqPath, path) || strings.Contains(refPath, path) {
			return true
		}
	}

	return false
}
