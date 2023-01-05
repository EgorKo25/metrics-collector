package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
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

	src = []byte(fmt.Sprintf("%s:%s:%f", metric.ID, metric.MType, *metric.Value))

	if h.Key == nil {
		return
	}

	hm := hmac.New(sha256.New, h.Key)
	hm.Write(src)
	hash = fmt.Sprintf("%x", hm.Sum(nil))

	log.Println("Вычисленный ", hash)
	log.Println("Имеющийся ", metric.Hash)
	if metric.Hash != "" && !hmac.Equal([]byte(metric.Hash), []byte(hash)) {
		return "", fmt.Errorf("not equal hash")
	}

	return
}
