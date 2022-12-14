package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
)

func main() {

	cfg := config.NewAgentConfig(2, 10)

	monitor := agent.NewMonitor(cfg)

	monitor.Run()
}
