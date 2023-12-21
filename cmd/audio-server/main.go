package main

import (
	"flag"

	"github.com/renan-campos/audio-server/pkg/metrics"
	"github.com/renan-campos/audio-server/pkg/storage"
	"github.com/renan-campos/audio-server/pkg/webserver"
)

func main() {
	listenAddr := flag.String("listen", "127.0.0.1", "The address to run the server on")
	port := flag.Int("port", 1323, "Port audio server will run on")
	runMetrics := flag.Bool("run-metrics", false, "Run the metrics server")

	flag.Parse()

	if *runMetrics {
		metricsServer := metrics.NewMetricsServer()
		go metricsServer.Run()
	}
	audioStorageService := storage.NewMemAudioStorageService()
	webserver := webserver.NewEchoWebServer(
		webserver.Parameters{Port: *port, ListenAddr: *listenAddr},
		webserver.Services{AudioStorage: audioStorageService},
	)
	webserver.Run()
}
