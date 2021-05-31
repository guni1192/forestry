package main

import (
	"github.com/guni1192/forestry/pkg/api"
	lokiV1 "github.com/guni1192/forestry/pkg/api/loki/v1"
	"github.com/guni1192/forestry/pkg/config"
)

func main() {
	s := api.NewServer(true)
	cfg := config.NewConfig()
	logStreamer := lokiV1.NewLogStreamer(lokiV1.ForestryLokiClient{LokiHost: cfg.LokiHost})

	s.GET("/health", api.Health)
	s.POST("/api/loki/v1/push", logStreamer.PushHandler)

	s.Logger.Fatal(s.Start(":1192"))
}
