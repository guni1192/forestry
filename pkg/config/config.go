package config

type Config struct {
	LokiHost string
}

func NewConfig() *Config {
	return &Config{LokiHost: "http://loki:3100"}
}
