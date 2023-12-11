package webserver

import (
	"fmt"
	"net/http"

	"github.com/renan-campos/audio-server/pkg/storage"

	"github.com/labstack/echo/v4"
)

var v0 struct {
	AdminRoutes   func(audioStorageService storage.AudioStorageService) EchoRoute
	RootEndpoints func(audioStorageService storage.AudioStorageService) []EchoEndpoint
} = struct {
	AdminRoutes   func(audioStorageService storage.AudioStorageService) EchoRoute
	RootEndpoints func(audioStorageService storage.AudioStorageService) []EchoEndpoint
}{
	AdminRoutes: func(audioStorageService storage.AudioStorageService) EchoRoute {
		return EchoRoute{}
	},
	RootEndpoints: func(audioStorageService storage.AudioStorageService) []EchoEndpoint {
		return []EchoEndpoint{
			{
				Path:   "/audio",
				Method: MethodGet,
				Handler: func(c echo.Context) error {
					listOfAudioIds, err := audioStorageService.ListAudio()
					if err != nil {
						return err
					}
					return c.JSON(http.StatusOK, listOfAudioIds)
				},
			},
			{
				Path:   "/audio/:id",
				Method: MethodGet,
				Handler: func(c echo.Context) error {
					id := c.Param("id")
					audioMetadata, err := audioStorageService.ListAudioMetadata(id)
					if err != nil {
						return err
					}
					return c.JSON(http.StatusOK, audioMetadata)
				},
			},
			{
				Path:   "/audio/:id/ogg",
				Method: MethodGet,
				Handler: func(c echo.Context) error {
					id := c.Param("id")
					// Specify the path to your Ogg sound file
					audioFilePath, err := audioStorageService.GetAudioFile(id)
					if err != nil {
						return err
					}

					// Set the appropriate headers for the HTTP response
					c.Response().Header().Set(echo.HeaderContentType, "audio/ogg")
					c.Response().Header().Set(echo.HeaderContentDisposition,
						fmt.Sprintf("attachment; filename=\"%s.ogg\"", id))

					// Serve the Ogg sound file as an HTTP response
					return c.File(string(audioFilePath))
				},
			},
		}
	},
}
