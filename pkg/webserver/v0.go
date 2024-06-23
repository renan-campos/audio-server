package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/renan-campos/audio-server/pkg/auth"
	"github.com/renan-campos/audio-server/pkg/handlers"
	"github.com/renan-campos/audio-server/pkg/storage"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/google/uuid"
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
					Path:   "/audio/:id/webm",
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
						return c.String(http.StatusOK, fmt.Sprintf("Uploaded webm file for %q\n", id))
					},
				},
			},
			Middlewares: []echo.MiddlewareFunc{
				/*
					middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
						// Be careful to use constant time comparison to prevent timing attacks
						// TODO: Don't hardcode the password!
						if subtle.ConstantTimeCompare([]byte(username), []byte("rcampos")) == 1 &&
							subtle.ConstantTimeCompare([]byte(password), []byte("relax")) == 1 {
							return true, nil
						}
						return false, nil
					}),
				*/
				auth.JwtAuth("auth.rcampos.net"),
			},
		}
	},
	RootEndpoints: func(audioStorageService storage.AudioStorageService) []EchoEndpoint {
		audioHandler := handlers.NewAudioHandler(handlers.AudioHandlerSpec{
			Storage: audioStorageService,
		})
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
				Path:    "/audio/:id/webm",
				Method:  MethodGet,
				Handler: audioHandler.GetAudioFile,
			},
			{
				Path:    "/audio/:id/webm",
				Method:  MethodHead,
				Handler: audioHandler.GetAudioHead,
			},
			{
				Path:   "/token",
				Method: MethodPost,
				Handler: func(c echo.Context) error {
					otp, err := io.ReadAll(c.Request().Body)
					if err != nil {
						return err
					}
					resp, err := http.Post(
						// Todo: auth endpoint
						fmt.Sprintf("%s/token", "http://auth.rcampos.net"),
						"application/json",
						bytes.NewBuffer(otp))
					if err != nil {
						return err
					}

					/* jwt logic { */
					defer resp.Body.Close()

					// Read response body
					rawJwt, err := io.ReadAll(resp.Body)
					if err != nil {
						return fmt.Errorf("Failed to read JWT: %v", err)
					}
					webToken, err := jwt.ParseSigned(string(rawJwt))
					if err != nil {
						return fmt.Errorf("Failed to parse JWT: %v", err)
					}
					var verifiedClaims jwt.Claims
					// Todo auth endpoint
					resp, err = http.Get("http://auth.rcampos.net/jwks")
					if err != nil {
						return fmt.Errorf("Http request failed: %v", err)
					}
					defer resp.Body.Close()

					var jwks jose.JSONWebKeySet
					marshalledJwks, err := io.ReadAll(resp.Body)
					if err != nil {
						return fmt.Errorf("Failed to read GET response body")
					}
					if err := json.Unmarshal(marshalledJwks, &jwks); err != nil {
						return fmt.Errorf("Failed to unmarshall jwks")
					}

					err = webToken.Claims(jwks.Keys[0], &verifiedClaims)
					if err != nil {
						log.Printf("Failed to verify jwt: %v", err)
						return err
					}
					log.Println("JWT verified!")
					err = verifiedClaims.Validate(jwt.Expected{
						Issuer: "Authentication-Server",
						Time:   time.Now(),
					})
					if err != nil {
						return fmt.Errorf("Failed to validate claims: %v", err)
					}
					log.Println("Claims validated sucessful.")

					/* jwt logic } */
					_, err = c.Response().Write(rawJwt)
					return err
				},
			},
		}
	},
}
