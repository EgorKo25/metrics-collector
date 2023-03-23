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
	"fmt"
	"log"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"

	_ "net/http/pprof"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {

	fmt.Printf("Build version: %s\nBuild date:    %s\nBuild commit:  %s\n", buildVersion, buildDate, buildCommit)

	cfg, err := config.NewAgentConfig()
	if err != nil {
		log.Fatalf("%s: %s", config.ErrFlagParse, err)
	}

	hsr := hashing.NewHash(cfg.Key)

	monitor := agent.NewMonitor(cfg, hsr)

	monitor.Run()
}
