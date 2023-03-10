package agent

import (
	"math/rand"
	"runtime"
	"testing"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"

	mems "github.com/shirou/gopsutil/v3/mem"
)

func BenchmarkMonitor_SendData(b *testing.B) {

	cfg := config.NewAgentConfig()

	hsr := hashing.NewHash(cfg.Key)

	m := NewMonitor(cfg, hsr)

	var mem runtime.MemStats
	var stats mems.VirtualMemoryStat
	var cpuInfo []float64

	m.RunMemStatListener(&mem)
	m.RunVirtMemCpuListener(&stats, &cpuInfo)

	for i := 0; i < b.N; i++ {
		b.Run("message sender", func(b *testing.B) {
			m.SendData(storage.Gauge(rand.Float64()), "RandomValue", "gauge")
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
		})
	}
}
