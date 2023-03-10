// Package agent пакет содержащий функционал агента
//
// Монитор - динамически проверяет состочние процессорных метрик и метрик памяти
// с заданным интервалом отправляет их на сервер
package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	mems "github.com/shirou/gopsutil/v3/mem"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

// Monitor структура монитора
type Monitor struct {
	config *config.ConfigurationAgent
	hasher *hashing.Hash

	pollCount storage.Counter
}

// NewMonitor конструтор структруы монитор
func NewMonitor(cfg *config.ConfigurationAgent, hsr *hashing.Hash) *Monitor {
	return &Monitor{
		config: cfg,
		hasher: hsr,
	}

}

// SendData отправляет метрики на сервер
func (m *Monitor) SendData(value storage.Gauge, name, mtype string) {
	var metric storage.Metric

	metric.ID = name
	metric.MType = mtype

	switch mtype {

	case "counter":
		tmp := storage.Counter(value)
		metric.Delta = &tmp

	case "gauge":
		metric.Value = &value
	}

	metric.Hash, _ = m.hasher.Run(&metric)

	dataJSON, err := json.Marshal(metric)
	if err != nil {
		log.Printf("Somethings went wrong: %s", err)
	}

	URL, _ := url.JoinPath("http://", m.config.Address, "update/")

	_, err = http.Post(URL, "application/json", bytes.NewBuffer(dataJSON))
	if err != nil {
		log.Printf("Somethings went wrong: %s", err)
	}
}

// RunMemStatListener считывает метрики памяти
func (m *Monitor) RunMemStatListener(mem *runtime.MemStats) {
	runtime.ReadMemStats(mem)
	m.pollCount++
}

// RunVirtMemCpuListener считывает метрики процессора
func (m *Monitor) RunVirtMemCpuListener(stats *mems.VirtualMemoryStat, cpuInfo *[]float64) {

	stats, _ = mems.VirtualMemory()
	*cpuInfo, _ = cpu.Percent(0, false)
	m.pollCount++
}

// Run запускает режим мониторинга в нескольких горутинах
func (m *Monitor) Run() {
	var mem runtime.MemStats
	var stats mems.VirtualMemoryStat
	var cpuInfo []float64

	tickerPoll := time.NewTicker(m.config.PollInterval)
	tickerReport := time.NewTicker(m.config.ReportInterval)

	for {
		select {

		case <-tickerPoll.C:
			go m.RunMemStatListener(&mem)
			go m.RunVirtMemCpuListener(&stats, &cpuInfo)

		case <-tickerReport.C:
			m.SendData(storage.Gauge(m.pollCount), "PollCount", "counter")
			m.SendData(storage.Gauge(rand.Float64()), "RandomValue", "gauge")
			m.SendData(storage.Gauge(mem.Alloc), "Alloc", "gauge")
			m.SendData(storage.Gauge(mem.BuckHashSys), "BuckHashSys", "gauge")
			m.SendData(storage.Gauge(mem.Frees), "Frees", "gauge")
			m.SendData(storage.Gauge(mem.GCCPUFraction), "GCCPUFraction", "gauge")
			m.SendData(storage.Gauge(mem.GCSys), "GCSys", "gauge")
			m.SendData(storage.Gauge(mem.HeapAlloc), "HeapAlloc", "gauge")
			m.SendData(storage.Gauge(mem.HeapIdle), "HeapIdle", "gauge")
			m.SendData(storage.Gauge(mem.HeapInuse), "HeapInuse", "gauge")
			m.SendData(storage.Gauge(mem.HeapObjects), "HeapObjects", "gauge")
			m.SendData(storage.Gauge(mem.HeapReleased), "HeapReleased", "gauge")
			m.SendData(storage.Gauge(mem.HeapSys), "HeapSys", "gauge")
			m.SendData(storage.Gauge(mem.LastGC), "LastGC", "gauge")
			m.SendData(storage.Gauge(mem.Lookups), "Lookups", "gauge")
			m.SendData(storage.Gauge(mem.MCacheInuse), "MCacheInuse", "gauge")
			m.SendData(storage.Gauge(mem.MCacheSys), "MCacheSys", "gauge")
			m.SendData(storage.Gauge(mem.MSpanInuse), "MSpanInuse", "gauge")
			m.SendData(storage.Gauge(mem.MSpanSys), "MSpanSys", "gauge")
			m.SendData(storage.Gauge(mem.Mallocs), "Mallocs", "gauge")
			m.SendData(storage.Gauge(mem.NextGC), "NextGC", "gauge")
			m.SendData(storage.Gauge(mem.NumForcedGC), "NumForcedGC", "gauge")
			m.SendData(storage.Gauge(mem.NumGC), "NumGC", "gauge")
			m.SendData(storage.Gauge(mem.OtherSys), "OtherSys", "gauge")
			m.SendData(storage.Gauge(mem.PauseTotalNs), "PauseTotalNs", "gauge")
			m.SendData(storage.Gauge(mem.StackInuse), "StackInuse", "gauge")
			m.SendData(storage.Gauge(mem.StackSys), "StackSys", "gauge")
			m.SendData(storage.Gauge(mem.Sys), "Sys", "gauge")
			m.SendData(storage.Gauge(mem.TotalAlloc), "TotalAlloc", "gauge")
			m.SendData(storage.Gauge(stats.Total), "TotalMemory", "gauge")
			m.SendData(storage.Gauge(stats.Free), "FreeMemory", "gauge")
			m.SendData(storage.Gauge(cpuInfo[0]), "CPUutilization1", "gauge")
		}
	}
}
