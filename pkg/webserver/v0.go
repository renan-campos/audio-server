package webserver

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/renan-campos/audio-server/pkg/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var v0 struct {
	AdminRoutes   func(audioStorageService storage.AudioStorageService) EchoRoute
	RootEndpoints func(audioStorageService storage.AudioStorageService) []EchoEndpoint
} = struct {
	AdminRoutes   func(audioStorageService storage.AudioStorageService) EchoRoute
	RootEndpoints func(audioStorageService storage.AudioStorageService) []EchoEndpoint
}{
	AdminRoutes: func(audioStorageService storage.AudioStorageService) EchoRoute {
		return EchoRoute{
			GroupPath: "/admin",
			Endpoints: []EchoEndpoint{
				{
					Path:   "/audio",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						newUuid := uuid.New().String()
						if err := audioStorageService.CreateEntry(newUuid); err != nil {
							return err
						}
						return c.String(http.StatusOK, newUuid)
					},
				},
				{
					Path:   "/audio/:id",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						id := c.Param("id")
						name := c.FormValue("name")
						audioStorageService.UpdateMetadata(id, storage.AudioMetadata{
							Name: name,
						})
						return c.String(http.StatusOK,
							fmt.Sprintf("Uploaded metadata for %q:\n{\n\tname: %q\n}\n", id, name))
					},
				},
				{
					Path:   "/audio/:id/ogg",
					Method: MethodPost,
					Handler: func(c echo.Context) error {
						// Get the uploaded file from the request
						file, err := c.FormFile("audioFile")
						if err != nil {
							fmt.Printf("error getting uploaded file: %v\n", err)
							return err
						}
						id := c.Param("id")
						if err := audioStorageService.UploadAudio(id, file); err != nil {
							fmt.Printf("error uploading audio: %v\n", err)
							fmt.Println(err)
							return err
						}
						return c.String(http.StatusOK, fmt.Sprintf("Uploaded ogg file for %q\n", id))
					},
				},
			},
			Middlewares: []echo.MiddlewareFunc{
				middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
					// Be careful to use constant time comparison to prevent timing attacks
					// TODO: Don't hardcode the password!
					if subtle.ConstantTimeCompare([]byte(username), []byte("rcampos")) == 1 &&
						subtle.ConstantTimeCompare([]byte(password), []byte("relax")) == 1 {
						return true, nil
					}
					return false, nil
				}),
			},
		}
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
