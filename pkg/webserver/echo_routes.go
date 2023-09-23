package webserver

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TODO: Eventually this function should take a set of handlers
func newEchoRoutes(e *echo.Echo) []EchoRoute {
	var msg string
	var secureMsg string

	return []EchoRoute{
		{
			GroupPath: "/",
			Endpoints: []EchoEndpoint{
				{
					Path:   "/",
					Method: MethodGet,
					Handler: func(c echo.Context) error {
						return c.String(http.StatusOK, fmt.Sprintf("%s\n", msg))
					},
				},
				{
					Path:   "/message",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						msg = c.FormValue("message")
						e.Logger.Infof("Message modified to %q", msg)
						return c.String(http.StatusOK, "Message modified\n")
					},
				},
			},
			Middlewares: []echo.MiddlewareFunc{},
		},
		{
			GroupPath: "/admin",
			Endpoints: []EchoEndpoint{
				{
					Path:   "/",
					Method: MethodGet,
					Handler: func(c echo.Context) error {
						return c.String(http.StatusOK, fmt.Sprintf("%s\n", secureMsg))
					},
				},
				{
					Path:   "/message",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						secureMsg = c.FormValue("message")
						e.Logger.Infof("Secure message modified to %q", msg)
						return c.String(http.StatusOK, "Message modified\n")
					},
				},
			},
			Middlewares: []echo.MiddlewareFunc{
				middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
					// Be careful to use constant time comparison to prevent timing attacks
					if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
						subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
						return true, nil
					}
					return false, nil
				}),
			},
		},
	}
}
