package agent

import (
	"bytes"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

func sendData(value storage.Gauge, name string, serializer *serializer.Serialize) {

	serializer.ID = name
	serializer.MType = "gauge"
	serializer.Value = &value

	dataJSON, err := serializer.Run()

	if err != nil {
		log.Fatalf("Somethings went wrong: %s", err)
	}

	URL, _ := url.JoinPath("http://127.0.0.1:8080", "update")

	_, err = http.Post(URL, "application/json", bytes.NewBuffer(dataJSON))
	if err != nil {
		log.Fatalf("Somethings went wrong: %s", err)
	}
}

type Monitor struct {
	serializer *serializer.Serialize

	pollCount storage.Counter

	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewMonitor(cfg *config.ConfigurationAgent, srl *serializer.Serialize) *Monitor {
	var mon Monitor
	mon.pollInterval = cfg.PollInterval
	mon.reportInterval = cfg.ReportInterval
	mon.serializer = srl
	return &mon

}
func (m *Monitor) Run() {
	var mem runtime.MemStats

	tickerPoll := time.NewTicker(m.pollInterval * time.Second)
	tickerReport := time.NewTicker(m.reportInterval * time.Second)

	for {
		select {
		case <-tickerPoll.C:
			runtime.ReadMemStats(&mem)
			m.pollCount++

		case <-tickerReport.C:
			sendData(storage.Gauge(mem.Alloc), "Alloc", m.serializer)
			sendData(storage.Gauge(mem.BuckHashSys), "BuckHashSys", m.serializer)
			sendData(storage.Gauge(mem.Frees), "Frees", m.serializer)
			sendData(storage.Gauge(mem.GCCPUFraction), "GCCPUFraction", m.serializer)
			sendData(storage.Gauge(mem.GCSys), "GCSys", m.serializer)
			sendData(storage.Gauge(mem.HeapAlloc), "HeapAlloc", m.serializer)
			sendData(storage.Gauge(mem.HeapIdle), "HeapIdle", m.serializer)
			sendData(storage.Gauge(mem.HeapInuse), "HeapInuse", m.serializer)
			sendData(storage.Gauge(mem.HeapObjects), "HeapObjects", m.serializer)
			sendData(storage.Gauge(mem.HeapReleased), "HeapReleased", m.serializer)
			sendData(storage.Gauge(mem.HeapSys), "HeapSys", m.serializer)
			sendData(storage.Gauge(mem.LastGC), "LastGC", m.serializer)
			sendData(storage.Gauge(mem.Lookups), "Lookups", m.serializer)
			sendData(storage.Gauge(mem.MCacheInuse), "MCacheInuse", m.serializer)
			sendData(storage.Gauge(mem.MCacheSys), "MCacheSys", m.serializer)
			sendData(storage.Gauge(mem.MSpanInuse), "MSpanInuse", m.serializer)
			sendData(storage.Gauge(mem.MSpanSys), "MSpanSys", m.serializer)
			sendData(storage.Gauge(mem.Mallocs), "Mallocs", m.serializer)
			sendData(storage.Gauge(mem.NextGC), "NextGC", m.serializer)
			sendData(storage.Gauge(mem.NumForcedGC), "NumForcedGC", m.serializer)
			sendData(storage.Gauge(mem.NumGC), "NumGC", m.serializer)
			sendData(storage.Gauge(mem.OtherSys), "OtherSys", m.serializer)
			sendData(storage.Gauge(mem.PauseTotalNs), "PauseTotalNs", m.serializer)
			sendData(storage.Gauge(mem.StackInuse), "StackInuse", m.serializer)
			sendData(storage.Gauge(mem.StackSys), "StackSys", m.serializer)
			sendData(storage.Gauge(mem.Sys), "Sys", m.serializer)
			sendData(storage.Gauge(mem.TotalAlloc), "TotalAlloc", m.serializer)
			sendData(storage.Gauge(rand.Float64()), "RandomValue", m.serializer)
		}
	}
}
