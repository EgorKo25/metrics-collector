package agent

import (
	"fmt"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

func sendData(value storage.Gauge, name string) {
	URL, err := url.JoinPath("http://127.0.0.1:8080", "update", "gauge", name, fmt.Sprintf("%f", value))
	if err != nil {
		log.Fatalf("Somethings went wrong: %s", err)
	}

	_, err = http.Post(URL, "text/plain", nil)
	if err != nil {
		log.Fatalf("Somethings went wrong: %s", err)
	}
}

type Monitor struct {
	pollCount storage.Counter

	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewMonitor(cfg *config.ConfigurationAgent) *Monitor {
	var mon Monitor
	mon.pollInterval = cfg.PollInterval
	mon.reportInterval = cfg.ReportInterval
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
			sendData(storage.Gauge(mem.Alloc), "Alloc")
			sendData(storage.Gauge(mem.BuckHashSys), "BuckHashSys")
			sendData(storage.Gauge(mem.Frees), "Frees")
			sendData(storage.Gauge(mem.GCCPUFraction), "GCCPUFraction")
			sendData(storage.Gauge(mem.GCSys), "GCSys")
			sendData(storage.Gauge(mem.HeapAlloc), "HeapAlloc")
			sendData(storage.Gauge(mem.HeapIdle), "HeapIdle")
			sendData(storage.Gauge(mem.HeapInuse), "HeapInuse")
			sendData(storage.Gauge(mem.HeapObjects), "HeapObjects")
			sendData(storage.Gauge(mem.HeapReleased), "HeapReleased")
			sendData(storage.Gauge(mem.HeapSys), "HeapSys")
			sendData(storage.Gauge(mem.LastGC), "LastGC")
			sendData(storage.Gauge(mem.Lookups), "Lookups")
			sendData(storage.Gauge(mem.MCacheInuse), "MCacheInuse")
			sendData(storage.Gauge(mem.MCacheSys), "MCacheSys")
			sendData(storage.Gauge(mem.MSpanInuse), "MSpanInuse")
			sendData(storage.Gauge(mem.MSpanSys), "MSpanSys")
			sendData(storage.Gauge(mem.Mallocs), "Mallocs")
			sendData(storage.Gauge(mem.NextGC), "NextGC")
			sendData(storage.Gauge(mem.NumForcedGC), "NumForcedGC")
			sendData(storage.Gauge(mem.NumGC), "NumGC")
			sendData(storage.Gauge(mem.OtherSys), "OtherSys")
			sendData(storage.Gauge(mem.PauseTotalNs), "PauseTotalNs")
			sendData(storage.Gauge(mem.StackInuse), "StackInuse")
			sendData(storage.Gauge(mem.StackSys), "StackSys")
			sendData(storage.Gauge(mem.Sys), "Sys")
			sendData(storage.Gauge(mem.TotalAlloc), "TotalAlloc")
			sendData(storage.Gauge(rand.Float64()), "RandomValue")
		}
	}
}
