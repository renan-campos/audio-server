package webserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TODO: Eventually this function should take a set of handlers
func newEchoRoutes(e *echo.Echo) []EchoRoute {
	return []EchoRoute{
		{
			GroupPath: "/v0",
			Endpoints: v0.RootEndpoints(),
			Middlewares: []echo.MiddlewareFunc{
				middleware.Static("./static"),
			},
			ChildRoutes: []EchoRoute{
				v0.AdminRoutes(),
			},
		},
	}
}
