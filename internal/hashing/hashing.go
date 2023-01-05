package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"

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

	switch metric.MType {
	case "gauge":
		src = []byte(fmt.Sprintf("%s:%s:%f", metric.ID, metric.MType, *metric.Value))
	case "counter":
		fmt.Println(metric.ID, metric.MType, *metric.Delta)
		src = []byte(fmt.Sprintf("%s:%s:%d", metric.ID, metric.MType, *metric.Delta))
	}

	if h.Key == nil {
		return "", nil
	}

	hm := hmac.New(sha256.New, h.Key)
	hm.Write(src)
	hash = fmt.Sprintf("%x", hm.Sum(nil))

	log.Println("Вычисленный ", hash)
	log.Println("Имеющийся ", metric.Hash)
	if metric.Hash != "" && !hmac.Equal([]byte(metric.Hash), []byte(hash)) {
		log.Println("not equal hash")
		return "", fmt.Errorf("not equal hash")
	}

	return hash, nil
}
