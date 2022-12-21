package serializer

import (
	"encoding/json"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

type Serialize struct {
	ID    string           `json:"id"`              // имя метрики
	MType string           `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *storage.Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *storage.Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewSerialize() *Serialize {
	return &Serialize{}
}

func (s Serialize) Run() ([]byte, error) {
	return json.Marshal(s)
}
