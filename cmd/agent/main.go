package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
)

func main() {
	//Формирую конфигурацию монитора
	cfg := config.NewAgentConfig(2, 10)
	//Создаю новый монитор
	monitor := agent.NewMonitor(cfg)
	//Запускаю монитор
	monitor.Run()
}
