package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
)

func main() {

	srl := serializer.NewSerialize()

	cfg := config.NewAgentConfig()

	monitor := agent.NewMonitor(cfg, srl)

	monitor.Run()
}
