package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
)

func main() {

	cfg := config.NewAgentConfig()

	monitor := agent.NewMonitor(cfg)

	monitor.Run()
}
