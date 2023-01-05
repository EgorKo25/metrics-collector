package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

type Hash struct {
	Key []byte
}

func MewHash(cfg string) *Hash {
	return &Hash{
		Key: []byte(cfg),
	}
}

func (h *Hash) Run(metric *storage.Metric) (hash string, err error) {

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
	hash = fmt.Sprintf("%b", hm.Sum(nil))

	if metric.Hash != "" {
		if hmac.Equal([]byte(metric.Hash), []byte(hash)) {
			return "", fmt.Errorf("not equal hash")
		}

	}

	return
}
