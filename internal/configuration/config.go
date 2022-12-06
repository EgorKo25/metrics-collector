package config

import "time"

type ConfigurationAgent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewAgentConfig(poll, report int) *ConfigurationAgent {
	var cfg ConfigurationAgent
	cfg.PollInterval = time.Duration(poll)
	cfg.ReportInterval = time.Duration(report)
	return &cfg
}
