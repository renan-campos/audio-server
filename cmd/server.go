package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	var msg string
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s\n", msg))
	})
	e.POST("/message", func(c echo.Context) error {
		msg = c.FormValue("message")
		e.Logger.Infof("Message modified to %q", msg)
		return c.String(http.StatusOK, "Message modified\n")
	})

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))
	var secureMsg string
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s\n", secureMsg))
	})
	g.POST("/message", func(c echo.Context) error {
		secureMsg = c.FormValue("message")
		e.Logger.Infof("Secure message modified to: %q", secureMsg)
		return c.String(http.StatusOK, "Secure message modified\n")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
