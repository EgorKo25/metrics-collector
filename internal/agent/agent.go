// Package agent пакет содержащий функционал агента
//
// Монитор - динамически проверяет состочние процессорных метрик и метрик памяти
// с заданным интервалом отправляет их на сервер
package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	mems "github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/encryption"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"github.com/EgorKo25/DevOps-Track-Yandex/proto/service"
)

var (
	ErrConstructor = errors.New("ошибка конструктора")
)

// Monitor структура монитора
type Monitor struct {
	config *config.ConfigurationAgent
	hasher *hashing.Hash
	enc    *encryption.Encryptor

	stats     *mems.VirtualMemoryStat
	PublicKey []byte

	pollCount storage.Counter

	shutdown bool
}

// NewMonitor конструтор структруы монитор
func NewMonitor(cfg *config.ConfigurationAgent, hsr *hashing.Hash, enc *encryption.Encryptor) (*Monitor, error) {

	return &Monitor{
		config: cfg,
		hasher: hsr,
		enc:    enc,
	}, nil

}

func (m *Monitor) SendGrpcData(value storage.Gauge, name, mtype string) {
	var metric service.Metric
	var me storage.Metric

	me.ID = name
	metric.Id = name

	if mtype == "gauge" {
		metric.Type = service.Metric_GAUGE
		me.MType = mtype
		me.Value = &value
	}
	if mtype == "gauge" {
		metric.Type = service.Metric_COUNTER
		me.MType = mtype
		delta := storage.Counter(value)
		me.Delta = &delta
	}

	metric.Hash, _ = m.hasher.Run(&me)

	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%s", err)
	}

	agent := service.NewServiceClient(conn)

	req, err := agent.TakeMetric(context.Background(), &service.MetricRequest{Metric: &metric})
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Println(req.Status)
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
		log.Printf("somethings went wrong: %s", err)
		return
	}

	dataJSON, err = m.enc.Encrypt(dataJSON)
	if err != nil {
		log.Printf("somethings went wrong: %s", err)
		return
	}

	URL, _ := url.JoinPath("http://", m.config.Address, "update/")

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(dataJSON))
	req.Header.Set("X-Real-IP", m.config.Address)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	defer resp.Body.Close()
}

// RunMemStatListener считывает метрики памяти
func (m *Monitor) RunMemStatListener(mem *runtime.MemStats) {
	runtime.ReadMemStats(mem)
	m.pollCount++
}

// RunVirtMemCpuListener считывает метрики процессора
func (m *Monitor) RunVirtMemCpuListener(cpuInfo *[]float64) {
	m.stats, _ = mems.VirtualMemory()
	*cpuInfo, _ = cpu.Percent(0, false)
	m.pollCount++

}

func (m *Monitor) GrpcMonitorRun() {
	var mem runtime.MemStats
	var cpuInfo []float64

	m.shutdown = true
	m.stats, _ = mems.VirtualMemory()

	tickerPoll := time.NewTicker(m.config.PollInterval)
	tickerReport := time.NewTicker(m.config.ReportInterval)

	for m.shutdown {

		select {

		case <-tickerPoll.C:
			go m.RunMemStatListener(&mem)
			go m.RunVirtMemCpuListener(&cpuInfo)

		case <-tickerReport.C:
			m.SendGrpcData(storage.Gauge(m.pollCount), "PollCount", "counter")
			m.SendGrpcData(storage.Gauge(rand.Float64()), "RandomValue", "gauge")
			m.SendGrpcData(storage.Gauge(mem.Alloc), "Alloc", "gauge")
			m.SendGrpcData(storage.Gauge(mem.BuckHashSys), "BuckHashSys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.Frees), "Frees", "gauge")
			m.SendGrpcData(storage.Gauge(mem.GCCPUFraction), "GCCPUFraction", "gauge")
			m.SendGrpcData(storage.Gauge(mem.GCSys), "GCSys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.HeapAlloc), "HeapAlloc", "gauge")
			m.SendGrpcData(storage.Gauge(mem.HeapIdle), "HeapIdle", "gauge")
			m.SendGrpcData(storage.Gauge(mem.HeapInuse), "HeapInuse", "gauge")
			m.SendGrpcData(storage.Gauge(mem.HeapObjects), "HeapObjects", "gauge")
			m.SendGrpcData(storage.Gauge(mem.HeapReleased), "HeapReleased", "gauge")
			m.SendGrpcData(storage.Gauge(mem.HeapSys), "HeapSys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.LastGC), "LastGC", "gauge")
			m.SendGrpcData(storage.Gauge(mem.Lookups), "Lookups", "gauge")
			m.SendGrpcData(storage.Gauge(mem.MCacheInuse), "MCacheInuse", "gauge")
			m.SendGrpcData(storage.Gauge(mem.MCacheSys), "MCacheSys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.MSpanInuse), "MSpanInuse", "gauge")
			m.SendGrpcData(storage.Gauge(mem.MSpanSys), "MSpanSys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.Mallocs), "Mallocs", "gauge")
			m.SendGrpcData(storage.Gauge(mem.NextGC), "NextGC", "gauge")
			m.SendGrpcData(storage.Gauge(mem.NumForcedGC), "NumForcedGC", "gauge")
			m.SendGrpcData(storage.Gauge(mem.NumGC), "NumGC", "gauge")
			m.SendGrpcData(storage.Gauge(mem.OtherSys), "OtherSys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.PauseTotalNs), "PauseTotalNs", "gauge")
			m.SendGrpcData(storage.Gauge(mem.StackInuse), "StackInuse", "gauge")
			m.SendGrpcData(storage.Gauge(mem.StackSys), "StackSys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.Sys), "Sys", "gauge")
			m.SendGrpcData(storage.Gauge(mem.TotalAlloc), "TotalAlloc", "gauge")
			m.SendGrpcData(storage.Gauge(m.stats.Total), "TotalMemory", "gauge")
			m.SendGrpcData(storage.Gauge(m.stats.Free), "FreeMemory", "gauge")
			m.SendGrpcData(storage.Gauge(cpuInfo[0]), "CPUutilization1", "gauge")
		}
	}
}

// Run запускает режим мониторинга в нескольких горутинах
func (m *Monitor) Run(sigs chan os.Signal) {
	var mem runtime.MemStats
	var cpuInfo []float64

	m.shutdown = true
	m.stats, _ = mems.VirtualMemory()

	tickerPoll := time.NewTicker(m.config.PollInterval)
	tickerReport := time.NewTicker(m.config.ReportInterval)

	for m.shutdown {

		select {

		case <-sigs:
			m.shutdown = false

		case <-tickerPoll.C:
			go m.RunMemStatListener(&mem)
			go m.RunVirtMemCpuListener(&cpuInfo)

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
			m.SendData(storage.Gauge(m.stats.Total), "TotalMemory", "gauge")
			m.SendData(storage.Gauge(m.stats.Free), "FreeMemory", "gauge")
			m.SendData(storage.Gauge(cpuInfo[0]), "CPUutilization1", "gauge")
		}
	}
}
