package config

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type ConfigurationAgent struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"5s"`
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func NewAgentConfig() *ConfigurationAgent {
	var cfg ConfigurationAgent
	_ = env.Parse(&cfg)
	return &cfg
}

type ConfigurationServer struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`

	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func NewServerConfig() *ConfigurationServer {
	var cfg ConfigurationServer
	_ = env.Parse(&cfg)
	return &cfg
}
