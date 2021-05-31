package main

import (
	"github.com/guni1192/forestry/pkg/api"
)

func main() {
	s := api.NewServer(true)

	s.GET("/health", api.Health)
	s.POST("/api/loki/v1/push", api.Health)

	s.Logger.Fatal(s.Start(":1192"))
}
