package webserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/renan-campos/audio-server/pkg/storage"
)

func NewEchoWebServer(audioStorageService storage.AudioStorageService) WebServer {
	e := &webServerImpl{
		Echo: echo.New(),
	}
	e.setupLogging()
	e.setupRoutes(audioStorageService)
	return e
}

type webServerImpl struct {
	*echo.Echo
}

func (e *webServerImpl) setupRoutes(audioStorageService storage.AudioStorageService) {
	routes := newEchoRoutes(e.Echo)
	for _, route := range routes {
		if route.GroupPath == "/" {
			e.handleRootEndpoints(route)
			continue
		}
		g := e.Group(route.GroupPath)
		for _, middleware := range route.Middlewares {
			g.Use(middleware)
		}
		for _, endpoint := range route.Endpoints {
			switch endpoint.Method {
			case MethodGet:
				g.GET(endpoint.Path, endpoint.Handler, endpoint.Middlewares...)
			case MethodPost:
				g.POST(endpoint.Path, endpoint.Handler, endpoint.Middlewares...)
			default:
				panic("Method not supported... yet.")
			}
		}
		for _, childRoute := range route.ChildRoutes {
			e.handleRoute(childRoute, g)
		}
	}
}

func (e *webServerImpl) handleRoute(route EchoRoute, parentGroup *echo.Group) {
	childGroup := parentGroup.Group(route.GroupPath)
	for _, middleware := range route.Middlewares {
		childGroup.Use(middleware)
	}
	for _, endpoint := range route.Endpoints {
		switch endpoint.Method {
		case MethodGet:
			childGroup.GET(endpoint.Path, endpoint.Handler, endpoint.Middlewares...)
		case MethodPost:
			childGroup.POST(endpoint.Path, endpoint.Handler, endpoint.Middlewares...)
		default:
			panic("Method not supported... yet.")
		}
	}
	for _, childRoute := range route.ChildRoutes {
		e.handleRoute(childRoute, childGroup)
	}
}

func (e *webServerImpl) handleRootEndpoints(route EchoRoute) {
	for _, middleware := range route.Middlewares {
		e.Use(middleware)
	}
	for _, endpoint := range route.Endpoints {
		switch endpoint.Method {
		case MethodGet:
			e.GET(endpoint.Path, endpoint.Handler)
		case MethodPost:
			e.POST(endpoint.Path, endpoint.Handler)
		default:
			panic("Method not supported... yet.")
		}
	}
}

func (e *webServerImpl) setupLogging() {
	e.Logger.SetLevel(log.INFO)
}

func (e *webServerImpl) Run() error {
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
	return nil
}
