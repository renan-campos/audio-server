package metrics

import (
	"errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type MetricsServer struct {
	app *echo.Echo
}

func NewMetricsServer() MetricsServer {
	metrics := echo.New()
	metrics.GET("/metrics", echoprometheus.NewHandler())
	return MetricsServer{
		app: metrics,
	}
}

func (s *MetricsServer) Run() {
	if err := s.app.Start(":1324"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
