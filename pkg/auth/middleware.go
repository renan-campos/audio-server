package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func JwtAuth(authUrl string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			jwt := jwtFromRequest(c.Request())
			if jwt == "" {
				c.Response().WriteHeader(http.StatusUnauthorized)
				return nil
			}
			if err := verifyJwt(jwt); err != nil {
				c.Response().WriteHeader(http.StatusUnauthorized)
				return nil
			}
			return next(c)
		}
	}
}

func jwtFromRequest(req *http.Request) string {
	header := req.Header.Get("Authorization")
	jwt := strings.Split(header, " ")
	if len(jwt) != 2 || jwt[0] != "Bearer" {
		return ""
	}
	return jwt[1]
}

func verifyJwt(token string) error {
	webToken, err := jwt.ParseSigned(string(token))
	if err != nil {
		return fmt.Errorf("Failed to parse JWT", err)
	}
	var verifiedClaims jwt.Claims
	// Todo auth endpoint
	resp, err := http.Get("http://auth.rcampos.net/jwks")
	if err != nil {
		fmt.Errorf("Http request failed:", err)
	}
	defer resp.Body.Close()

	var jwks jose.JSONWebKeySet
	marshalledJwks, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Failed to read GET response body")
	}
	if err := json.Unmarshal(marshalledJwks, &jwks); err != nil {
		fmt.Errorf("Failed to unmarshall jwks")
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
		fmt.Errorf("Failed to validate claims:", err)
	}
	log.Println("Claims validated sucessful.\n")

	/* jwt logic } */
	return nil
}
