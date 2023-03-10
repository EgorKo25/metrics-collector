package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"

	_ "net/http/pprof"
)

func main() {

	cfg := config.NewAgentConfig()

	hsr := hashing.NewHash(cfg.Key)

	monitor := agent.NewMonitor(cfg, hsr)

	monitor.Run()
}
