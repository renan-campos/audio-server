package main

import (
	"github.com/renan-campos/audio-server/pkg/metrics"
	"github.com/renan-campos/audio-server/pkg/storage"
	"github.com/renan-campos/audio-server/pkg/webserver"
)

func main() {
	metricsServer := metrics.NewMetricsServer()
	go metricsServer.Run()
	audioStorageService := storage.NewMemAudioStorageService()
	webserver := webserver.NewEchoWebServer(audioStorageService)
	webserver.Run()
}
