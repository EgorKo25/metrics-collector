// Пакет с логикой агента
//
// Приложение собирается командой:
//
//	go build agent
//
// Запускается:
//
//	./agent
//
// Или
//
//	go run main.go
package main

import (
	"log"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"

	_ "net/http/pprof"
)

func main() {

	cfg, err := config.NewAgentConfig()
	if err != nil {
		log.Fatalf("%s: %s", config.ErrFlagParse, err)
		return
	}

	hsr := hashing.NewHash(cfg.Key)

	monitor := agent.NewMonitor(cfg, hsr)

	monitor.Run()
}
