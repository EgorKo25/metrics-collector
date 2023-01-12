package file

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

type Save struct {
	cfg  *config.ConfigurationServer
	strg *storage.MetricStorage

	file   *os.File
	writer *bufio.Writer
}

func NewSave(cfg *config.ConfigurationServer, strg *storage.MetricStorage) *Save {
	file, _ := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	return &Save{
		cfg:    cfg,
		strg:   strg,
		file:   file,
		writer: bufio.NewWriter(file),
	}
}

func (s *Save) Close() error {
	return s.file.Close()
}

func (s *Save) WriteAll() (err error) {
	var metric storage.Metric
	var data []byte

	for k, v := range s.strg.Metrics {

		metric.ID = k
		metric.MType = v.MType

		if v.MType == "gauge" {
			metric.Value = v.Value
			metric.Delta = nil
		}

		if v.MType == "counter" {
			metric.Value = nil
			metric.Delta = v.Delta
		}

		if data, err = json.Marshal(metric); err != nil {
			return
		}
		if _, err = s.writer.Write(data); err != nil {
			return
		}
		if err = s.writer.WriteByte('\n'); err != nil {
			return
		}
	}

	return s.writer.Flush()
}

func (s *Save) Run() error {
	tickerSave := time.NewTicker(s.cfg.StoreInterval)

	for {
		select {
		case <-tickerSave.C:
			if err := s.WriteAll(); err != nil {
				return err
			}
		}
	}
}

type Read struct {
	cfg  *config.ConfigurationServer
	strg *storage.MetricStorage

	file   *os.File
	reader *bufio.Reader
}

func NewRead(cfg *config.ConfigurationServer, strg *storage.MetricStorage) (*Read, error) {

	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)

	reader := &Read{
		cfg:    cfg,
		strg:   strg,
		file:   file,
		reader: bufio.NewReader(file),
	}

	if cfg.Restore {
		err = reader.readAll()
		if err != nil {
			log.Println("file read error: ", err)
		}
	}

	return reader, err
}

func (r *Read) readAll() (err error) {

	var data []byte
	var metric storage.Metric

	for {

		if data, err = r.reader.ReadBytes('\n'); err != nil {
			return err
		}

		if err = json.Unmarshal(data, &metric); err != nil {
			return err
		}

		if metric.MType == "gauge" {
			r.strg.Metrics[metric.ID] = metric
		}

		if metric.MType == "counter" {
			r.strg.Metrics[metric.ID] = metric

		}

	}

}
