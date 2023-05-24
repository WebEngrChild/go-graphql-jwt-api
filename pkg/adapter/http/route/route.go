package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"

	"golang.org/x/xerrors"

	"github.com/WebEngrChild/go-graphql-server/pkg/adapter/http/handler"
	authMiddleware "github.com/WebEngrChild/go-graphql-server/pkg/adapter/http/middleware"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/config"
	"github.com/WebEngrChild/go-graphql-server/pkg/lib/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Route interface {
	InitRouting(*config.Config) (*echo.Echo, error)
}

type InitRoute struct {
	Ch handler.Csrf
	Lh handler.Login
	Gh handler.Graph
	Ph http.HandlerFunc
	Am authMiddleware.Auth
}

func NewInitRoute(ch handler.Csrf, lh handler.Login, gh handler.Graph, ph http.HandlerFunc, am authMiddleware.Auth) Route {
	InitRoute := InitRoute{ch, lh, gh, ph, am}
	return &InitRoute
}

func (i *InitRoute) InitRouting(cfg *config.Config) (*echo.Echo, error) {
	e := echo.New()

	cookieDomain := ""
	if cfg.Env == "prd" {
		cookieDomain = "." + cfg.AppDomain
	}

	// middleware
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{cfg.FrontURL},
			AllowCredentials: true,
		}),
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			CookiePath:     "/",
			CookieSecure:   true,
			CookieDomain:   cookieDomain,
			CookieSameSite: http.SameSiteNoneMode,
			Skipper: func(c echo.Context) bool {
				if strings.Contains(c.Request().URL.Path, "/healthcheck") {
					return true
				}
				if strings.Contains(c.Request().URL.Path, "/playground") {
					return true
				}
				if strings.Contains(c.Request().URL.Path, "/query") {
					return true
				}
				return false
			},
		}),
		i.Am.AuthMiddleware,
	)

	// Validator
	e.Validator = validator.NewValidator()

	// Custom Error Handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Route
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "New deployment test")
	})
	e.GET("/csrf-cookie", i.Ch.CsrfHandler())
	e.POST("/login", i.Lh.LoginHandler())
	e.GET("/logout", i.Lh.LogoutHandler())
	e.POST("/query", i.Gh.QueryHandler())
	e.GET("/playground", func(c echo.Context) error {
		i.Ph.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Start Server
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		return nil, xerrors.Errorf("fail to start port:%s %w", cfg.Port, err)
	}

	return e, nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	sentry.CaptureException(fmt.Errorf("handler err: %w", err))

	c.Logger().Error(err)

	if err := c.JSON(http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	}); err != nil {
		c.Logger().Error(err)
	}
}
