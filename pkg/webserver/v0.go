package webserver

import (
	"net/http"

	"github.com/renan-campos/audio-server/pkg/storage"

	"github.com/labstack/echo/v4"
)

var v0 struct {
	RootEndpoints func(movieStorageService storage.MovieStorageService) []EchoEndpoint
} = struct {
	RootEndpoints func(movieStorageService storage.MovieStorageService) []EchoEndpoint
}{
	RootEndpoints: func(movieStorageService storage.MovieStorageService) []EchoEndpoint {
		return []EchoEndpoint{
			{
				Path:   "/movie",
				Method: MethodPost,
				Handler: func(c echo.Context) error {
					name := c.FormValue("name")
					if err := movieStorageService.CreateEntry(name); err != nil {
						return err
					}
					return c.String(http.StatusOK, name)
				},
			},
			{
				Path:   "/movies",
				Method: MethodGet,
				Handler: func(c echo.Context) error {
					movieList, err := movieStorageService.ListMovies()
					if err != nil {
						return err
					}
					return c.JSON(http.StatusOK, movieList)
				},
			},
		}
	},
}
