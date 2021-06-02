package main

import (
	"fmt"
	"os"

	"github.com/guni1192/forestry/pkg/api"
	lokiV1 "github.com/guni1192/forestry/pkg/api/loki/v1"
	"github.com/guni1192/forestry/pkg/config"
	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Use:   "forestry-server",
		Short: "Logging Pub/Sub Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			lokiHost, err := cmd.Flags().GetString("loki-host")
			if err != nil {
				return err
			}
			cfg := config.NewConfig()
			cfg.LokiHost = lokiHost
			run(*cfg)
			return nil
		},
	}
)

func init() {
	rootCommand.PersistentFlags().StringP("loki-host", "", "http://localhost:3100", "Grafana Loki hostname")
}

func run(cfg config.Config) {
	s := api.NewServer(true)
	s.Logger.Info(cfg)

	logStreamer := lokiV1.NewLogStreamer(lokiV1.ForestryLokiClient{LokiHost: cfg.LokiHost})

	s.GET("/health", api.Health)
	s.POST("/api/loki/v1/push", logStreamer.PushHandler)

	s.Logger.Fatal(s.Start(":1192"))
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
