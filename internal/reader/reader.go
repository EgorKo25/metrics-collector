package reader

import (
	"bufio"
	"encoding/json"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"io"
	"os"
)

type Read struct {
	cfg  *config.ConfigurationServer
	srl  *serializer.Serialize
	strg *storage.MetricStorage

	file   *os.File
	reader *bufio.Reader
}

func NewRead(cfg *config.ConfigurationServer, strg *storage.MetricStorage, srl *serializer.Serialize) (*Read, error) {
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
	return &Read{
		cfg:    cfg,
		strg:   strg,
		srl:    srl,
		file:   file,
		reader: bufio.NewReader(file),
	}, err
}

func (r *Read) ReadAll() (data []byte, err error) {
	r.srl.Clean()
	for {

		if data, err = r.reader.ReadBytes('\n'); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(data, r.srl); err != nil {
			return nil, err
		}
		if r.srl.MType == "gauge" {
			r.strg.MetricsGauge[r.srl.ID] = *r.srl.Value
			r.srl.Clean()
		}
		if r.srl.MType == "counter" {
			r.strg.MetricsCounter[r.srl.ID] = *r.srl.Delta
			r.srl.Clean()
		}
		_, isEnd := r.file.Read(data)
		if isEnd == io.EOF {
			return
		}

	}

	return
}
