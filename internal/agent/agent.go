package agent

import (
	"bytes"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"math/rand"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

type Monitor struct {
	serializer *serializer.Serialize
	config     *config.ConfigurationAgent

	pollCount storage.Counter
}

func NewMonitor(cfg *config.ConfigurationAgent, srl *serializer.Serialize) *Monitor {
	var mon Monitor
	mon.config = cfg
	mon.serializer = srl
	return &mon

}
func (m *Monitor) sendData(value storage.Gauge, name, Mtype string) {
	m.serializer.Clean()

	m.serializer.ID = name
	m.serializer.MType = Mtype

	switch Mtype {
	case "counter":
		tmp := storage.Counter(value)
		m.serializer.Delta = &tmp

	case "gauge":
		m.serializer.Value = &value
	}

	dataJSON, err := m.serializer.Run()
	if err != nil {
		log.Printf("Somethings went wrong: %s", err)
	}

	addr := m.config.Address
	URL, _ := url.JoinPath("http://", addr, "update/")

	_, err = http.Post(URL, "application/json", bytes.NewBuffer(dataJSON))
	if err != nil {
		log.Printf("Somethings went wrong: %s", err)
	}
}
func (m *Monitor) Run() {
	var mem runtime.MemStats

	tickerPoll := time.NewTicker(m.config.PollInterval)
	tickerReport := time.NewTicker(m.config.ReportInterval)

	for {
		select {
		case <-tickerPoll.C:
			runtime.ReadMemStats(&mem)
			m.pollCount++

		case <-tickerReport.C:
			m.sendData(storage.Gauge(m.pollCount), "PollCount", "counter")
			m.sendData(storage.Gauge(rand.Float64()), "RandomValue", "gauge")
			m.sendData(storage.Gauge(mem.Alloc), "Alloc", "gauge")
			m.sendData(storage.Gauge(mem.BuckHashSys), "BuckHashSys", "gauge")
			m.sendData(storage.Gauge(mem.Frees), "Frees", "gauge")
			m.sendData(storage.Gauge(mem.GCCPUFraction), "GCCPUFraction", "gauge")
			m.sendData(storage.Gauge(mem.GCSys), "GCSys", "gauge")
			m.sendData(storage.Gauge(mem.HeapAlloc), "HeapAlloc", "gauge")
			m.sendData(storage.Gauge(mem.HeapIdle), "HeapIdle", "gauge")
			m.sendData(storage.Gauge(mem.HeapInuse), "HeapInuse", "gauge")
			m.sendData(storage.Gauge(mem.HeapObjects), "HeapObjects", "gauge")
			m.sendData(storage.Gauge(mem.HeapReleased), "HeapReleased", "gauge")
			m.sendData(storage.Gauge(mem.HeapSys), "HeapSys", "gauge")
			m.sendData(storage.Gauge(mem.LastGC), "LastGC", "gauge")
			m.sendData(storage.Gauge(mem.Lookups), "Lookups", "gauge")
			m.sendData(storage.Gauge(mem.MCacheInuse), "MCacheInuse", "gauge")
			m.sendData(storage.Gauge(mem.MCacheSys), "MCacheSys", "gauge")
			m.sendData(storage.Gauge(mem.MSpanInuse), "MSpanInuse", "gauge")
			m.sendData(storage.Gauge(mem.MSpanSys), "MSpanSys", "gauge")
			m.sendData(storage.Gauge(mem.Mallocs), "Mallocs", "gauge")
			m.sendData(storage.Gauge(mem.NextGC), "NextGC", "gauge")
			m.sendData(storage.Gauge(mem.NumForcedGC), "NumForcedGC", "gauge")
			m.sendData(storage.Gauge(mem.NumGC), "NumGC", "gauge")
			m.sendData(storage.Gauge(mem.OtherSys), "OtherSys", "gauge")
			m.sendData(storage.Gauge(mem.PauseTotalNs), "PauseTotalNs", "gauge")
			m.sendData(storage.Gauge(mem.StackInuse), "StackInuse", "gauge")
			m.sendData(storage.Gauge(mem.StackSys), "StackSys", "gauge")
			m.sendData(storage.Gauge(mem.Sys), "Sys", "gauge")
			m.sendData(storage.Gauge(mem.TotalAlloc), "TotalAlloc", "gauge")

		}
	}
}
