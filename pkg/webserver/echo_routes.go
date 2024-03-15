package webserver

import (
	"github.com/renan-campos/audio-server/pkg/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newEchoRoutes(e *echo.Echo, movieStorageService storage.MovieStorageService) []EchoRoute {
	return []EchoRoute{
		{
			GroupPath: "/",
			Middlewares: []echo.MiddlewareFunc{
				middleware.Static(staticFileLocation),
			},
		},
		{
			GroupPath: "/v0",
			Endpoints: v0.RootEndpoints(movieStorageService),
			Middlewares: []echo.MiddlewareFunc{
				middleware.Logger(),
			},
			ChildRoutes: []EchoRoute{},
		},
	}
}
