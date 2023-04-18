// Package config
// пакет формирует конфигурации для агента и сервера
//
// Конфигурация агента создается при помощи команды:
//
//	agent := NewAgentConfig().
//
// Конфигурация сервера создается при помощи команды:
//
//	agent := NewServerConfig().
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

var (
	ErrFlagParse = errors.New("failed to get flag values")
)

// ConfigurationAgent структура конфигурации агента
type ConfigurationAgent struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	Address        string        `env:"ADDRESS" json:"address"`
	Key            string        `env:"KEY"`
	CryptoKey      string        `env:"CRYPTO_KEY" json:"crypto_key"`
	configFile     string        `env:"CONFIG"`
}

// NewAgentConfig конструтор конфигурации объекта
func NewAgentConfig() (*ConfigurationAgent, error) {
	var cfg ConfigurationAgent
	flag.StringVar(&cfg.Address,
		"a", "127.0.0.1:8080",
		"listening address of the server",
	)
	flag.StringVar(&cfg.Key,
		"k", "",
		"traffic encryption key",
	)
	flag.StringVar(&cfg.Key,
		"crypto-key", "",
		"traffic encryption key",
	)
	flag.StringVar(&cfg.configFile,
		"c", "",
		"file of configuration",
	)
	flag.DurationVar(&cfg.PollInterval,
		"p", time.Second*2,
		"timeout for update metrics",
	)
	flag.DurationVar(&cfg.ReportInterval,
		"r", time.Second*5,
		"timeout for report metrics",
	)

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.configFile != "" {
		file, _ := os.ReadFile(cfg.configFile)

		_ = json.Unmarshal(file, &cfg)
	}

	return &cfg, nil
}

// ConfigurationServer структура конфигурации агента
type ConfigurationServer struct {
	Address       string        `env:"ADDRESS" json:"address"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile     string        `env:"STORE_FILE" json:"store_file"`
	Restore       bool          `env:"RESTORE" json:"restore"`
	Key           string        `env:"KEY"`
	DB            string        `env:"DATABASE_DSN" json:"database_dsn"`
	CryptoKey     string        `env:"CRYPTO_KEY" json:"crypto_key"`
	TrustedSubnet string        `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	configFile    string        `env:"CONFIG"`
}

// NewServerConfig конструктор конфигурации объекта
func NewServerConfig() (*ConfigurationServer, error) {
	var cfg ConfigurationServer

	flag.StringVar(&cfg.Address,
		"a", "127.0.0.1:8080",
		"listening address of the server",
	)
	flag.StringVar(&cfg.Key,
		"k", "",
		"traffic encryption key",
	)
	flag.StringVar(&cfg.TrustedSubnet,
		"t", "",
		"trusted subnet",
	)
	flag.StringVar(&cfg.CryptoKey,
		"crypto-key", "",
		"traffic encryption key",
	)
	flag.StringVar(&cfg.configFile,
		"c", "",
		"file of configuration",
	)
	flag.StringVar(&cfg.DB,
		"d", "",
		"string with the address of the connection to the database",
	)
	flag.DurationVar(&cfg.StoreInterval,
		"i", time.Second*300,
		"timeout for save metrics",
	)
	flag.BoolVar(&cfg.Restore,
		"r", true,
		"recover files from storage",
	)
	flag.StringVar(&cfg.StoreFile,
		"f", "/tmp/devops-metrics-db.json",
		"file for saving metrics",
	)

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.configFile != "" {
		file, _ := os.ReadFile(cfg.configFile)

		_ = json.Unmarshal(file, &cfg)
	}

	return &cfg, nil
}
