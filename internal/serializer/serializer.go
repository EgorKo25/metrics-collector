package serializer

import (
	"encoding/json"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

type Serialize struct {
	ID    string           `json:"id"`
	MType string           `json:"type"`
	Delta *storage.Counter `json:"delta,omitempty"`
	Value *storage.Gauge   `json:"value,omitempty"`
}

func NewSerialize() *Serialize {
	return &Serialize{}
}

func (s Serialize) Run() ([]byte, error) {
	return json.Marshal(s)
}
