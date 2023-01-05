package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

type Hash struct {
	Key []byte
}

func MewHash(cfg *config.ConfigurationAgent) *Hash {
	return &Hash{
		Key: []byte(cfg.Key),
	}
}

func (h *Hash) Run(metric *storage.Metric) (hash string) {

	var src []byte

	if metric.MType == "gauge" {
		src = []byte(fmt.Sprintf("%s:guage:%f", metric.ID, *metric.Value))
	}
	if metric.MType == "counter" {
		src = []byte(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta))
	}

	if h.Key == nil {
		return
	}

	hm := hmac.New(sha256.New, h.Key)
	hm.Write(src)
	hash = string(hm.Sum(nil))
	return
}
