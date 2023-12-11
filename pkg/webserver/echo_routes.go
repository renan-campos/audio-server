package webserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newEchoRoutes(e *echo.Echo) []EchoRoute {
	return []EchoRoute{
		{
			GroupPath: "/v0",
			Middlewares: []echo.MiddlewareFunc{
				middleware.Static(staticFileLocation),
			},
		},
	}
}
