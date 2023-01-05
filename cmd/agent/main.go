package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
)

func main() {

	cfg := config.NewAgentConfig()

	hsr := hashing.MewHash(cfg)

	monitor := agent.NewMonitor(cfg, hsr)

	monitor.Run()
}
