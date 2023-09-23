package webserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func NewEchoWebServer() WebServer {
	e := &webServerImpl{
		Echo: echo.New(),
	}
	e.Logger.SetLevel(log.INFO)
	e.setupRoutes()
	return e
}

type webServerImpl struct {
	*echo.Echo
}

func (e *webServerImpl) setupRoutes() {
	routes := newEchoRoutes(e.Echo)
	for _, route := range routes {
		g := e.Group(route.GroupPath)
		for _, middleware := range route.Middlewares {
			g.Use(middleware)
		}
		for _, endpoint := range route.Endpoints {
			switch endpoint.Method {
			case MethodGet:
				g.GET(endpoint.Path, endpoint.Handler)
			case MethodPost:
				g.POST(endpoint.Path, endpoint.Handler)
			default:
				panic("Method not supported... yet.")
			}
		}
	}
}

func (e *webServerImpl) Run() error {
	e.Logger.Fatal(e.Start(":1323"))
	return nil
}
