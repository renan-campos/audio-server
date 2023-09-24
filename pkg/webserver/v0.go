package webserver

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var v0 struct {
	AdminRoutes   func() EchoRoute
	RootEndpoints func() []EchoEndpoint
} = struct {
	AdminRoutes   func() EchoRoute
	RootEndpoints func() []EchoEndpoint
}{
	AdminRoutes: func() EchoRoute {
		return EchoRoute{
			GroupPath: "/admin",
			Endpoints: []EchoEndpoint{
				{
					Path:   "/audio",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						newUUID := uuid.New()
						return c.String(http.StatusOK, fmt.Sprintf("Audio resource created: %q\n", newUUID.String()))
					},
				},
				{
					Path:   "/audio/:id",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						id := c.Param("id")
						secureName := c.FormValue("name")
						return c.String(http.StatusOK,
							fmt.Sprintf("Uploaded metadata for %q:\n{\n\tname: %q\n}\n", id, secureName))
					},
				},
				{
					Path:   "/audio/:id/ogg",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						id := c.Param("id")
						return c.String(http.StatusOK, fmt.Sprintf("Uploaded ogg file for %q\n", id))
					},
				},
			},
			Middlewares: []echo.MiddlewareFunc{
				middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
					// Be careful to use constant time comparison to prevent timing attacks
					// TODO: Don't hardcode the password!
					if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
						subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
						return true, nil
					}
					return false, nil
				}),
			},
		}
	},
	RootEndpoints: func() []EchoEndpoint {
		return []EchoEndpoint{
			{
				Path:   "/audio",
				Method: MethodGet,
				Handler: func(c echo.Context) error {
					return c.String(http.StatusOK, "List of audio metadata\n")
				},
			},
			{
				Path:   "/audio/:id",
				Method: MethodGet,
				Handler: func(c echo.Context) error {
					id := c.Param("id")
					return c.String(http.StatusOK, fmt.Sprintf("List of audio %q metadata\n", id))
				},
			},
			{
				Path:   "/audio/:id/ogg",
				Method: MethodGet,
				Handler: func(c echo.Context) error {
					id := c.Param("id")
					return c.String(http.StatusOK, fmt.Sprintf("ogg file of %q\n", id))
				},
			},
		}
	},
}
