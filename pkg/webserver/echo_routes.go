package webserver

import (
	"github.com/renan-campos/audio-server/pkg/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newEchoRoutes(e *echo.Echo, audioStorageService storage.AudioStorageService) []EchoRoute {
	return []EchoRoute{
		{
			GroupPath: "/",
			Middlewares: []echo.MiddlewareFunc{
				middleware.Static(staticFileLocation),
			},
		},
		{
			GroupPath: "/v0",
			Endpoints: v0.RootEndpoints(audioStorageService),
			Middlewares: []echo.MiddlewareFunc{
				middleware.Logger(),
			},
			ChildRoutes: []EchoRoute{
				v0.AdminRoutes(audioStorageService),
			},
		},
	}
}
