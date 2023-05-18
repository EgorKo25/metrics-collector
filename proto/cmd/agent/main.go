package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/agent"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/encryption"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
)

func main() {

	cfg, _ := config.NewAgentConfig()
	hsr := hashing.NewHash(cfg.Key)
	enc, _ := encryption.NewEncryptor(cfg.CryptoKey, "public")

	grpcAgent, _ := agent.NewMonitor(cfg, hsr, enc)

	grpcAgent.GrpcMonitorRun()
}
