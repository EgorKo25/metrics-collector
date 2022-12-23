package saver

import (
	"bufio"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"os"
	"time"
)

type Save struct {
	cfg  *config.ConfigurationServer
	srl  *serializer.Serialize
	strg *storage.MetricStorage

	file   *os.File
	writer *bufio.Writer
}

func NewSave(cfg *config.ConfigurationServer, strg *storage.MetricStorage, srl *serializer.Serialize) *Save {
	file, _ := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	return &Save{
		cfg:    cfg,
		strg:   strg,
		srl:    srl,
		file:   file,
		writer: bufio.NewWriter(file),
	}
}

func (s *Save) Close() error {
	return s.file.Close()
}

func (s *Save) WriteAll() (err error) {
	var data []byte
	s.srl.Clean()

	s.srl.MType = "gauge"
	for k, v := range s.strg.MetricsGauge {
		s.srl.ID = k
		s.srl.Value = &v

		if data, err = s.srl.Run(); err != nil {
			return
		}
		if _, err = s.writer.Write(data); err != nil {
			return
		}
		if err = s.writer.WriteByte('\n'); err != nil {
			return
		}
	}
	s.srl.MType = "counter"
	for k, v := range s.strg.MetricsCounter {
		s.srl.ID = k
		s.srl.Delta = &v

		if data, err = s.srl.Run(); err != nil {
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
