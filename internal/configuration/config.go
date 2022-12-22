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
