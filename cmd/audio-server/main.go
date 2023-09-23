package main

import (
	"github.com/renan-campos/audio-server/pkg/webserver"
)

func main() {
	webserver := webserver.NewEchoWebServer()
	webserver.Run()
}
