package webserver

import (
	"github.com/labstack/echo/v4"
)

type WebServer interface {
	Run() error
}

type EchoRoute struct {
	Endpoints   []EchoEndpoint
	GroupPath   string
	Middlewares []echo.MiddlewareFunc
	ChildRoutes []EchoRoute
}

type EchoEndpoint struct {
	Handler     echo.HandlerFunc
	Method      HttpMethod
	Middlewares []echo.MiddlewareFunc
	Path        string
}

// I'm suprised the standard library doesn't define such a type
type HttpMethod string

// This was copied straight from net/http
const (
	MethodGet     HttpMethod = "GET"
	MethodHead    HttpMethod = "HEAD"
	MethodPost    HttpMethod = "POST"
	MethodPut     HttpMethod = "PUT"
	MethodPatch   HttpMethod = "PATCH" // RFC 5789
	MethodDelete  HttpMethod = "DELETE"
	MethodConnect HttpMethod = "CONNECT"
	MethodOptions HttpMethod = "OPTIONS"
	MethodTrace   HttpMethod = "TRACE"
)
