package main

import (
	"log"
	"net/http"
	rt "runtime"
	"strconv"
	"time"
)

type guage uint64

type counter int

type config struct {
	pollInterval   time.Duration
	reportInterval time.Duration
}

func copy_map(metric_list *map[string]*uint64) map[string]uint64 {
	old_metric_list := make(map[string]uint64)
	for k, v := range *metric_list {
		old_metric_list[k] = *v
	}
	return old_metric_list
}

func CountAndCompare(m *rt.MemStats, metric_list *map[string]*uint64, PollCounter *counter) {
	old_metric_list := copy_map(metric_list)
	rt.ReadMemStats(m)
	for k, _ := range *metric_list {
		if old_metric_list[k] != *(*metric_list)[k] {
			*PollCounter++
		}
	}
}

func UpdateMetric(PollCounter *counter) (metric_list map[string]*uint64) {
	var m rt.MemStats
	rt.ReadMemStats(&m)

	GCCPUFraction := uint64(m.GCCPUFraction)
	metric_list = map[string]*uint64{
		"Alloc":         &m.Alloc,
		"BuckHashSys":   &m.BuckHashSys,
		"Frees":         &m.Frees,
		"GCCPUFraction": &GCCPUFraction,
		"GCSys":         &m.GCSys,
		"HeapAlloc":     &m.HeapAlloc,
		"HeapIdle":      &m.HeapIdle,
	}

	CountAndCompare(&m, &metric_list, PollCounter)
	return
}

func MetricSender(metrcs_list *map[string]*uint64) {
	for k, v := range *metrcs_list {

		url := "http://127.0.0.1:8080/update/" + "guage" + "/" + k + "/" + strconv.Itoa(int(*v)) + "/"
		_, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

	}

}

func Monitor() {
	var m rt.MemStats
	rt.ReadMemStats(&m)

	var PollCount counter = 0
	var metrics_list map[string]*uint64

	var dur_conf config
	dur_conf.pollInterval = 2
	dur_conf.reportInterval = 10

	tickerInterval := time.NewTicker(dur_conf.pollInterval * time.Second)
	tickerReport := time.NewTicker(dur_conf.reportInterval * time.Second)

	for {
		select {
		case <-tickerInterval.C:
			metrics_list = UpdateMetric(&PollCount)

		case <-tickerReport.C:
			MetricSender(&metrics_list)
		}
	}

}
func main() {
	Monitor()
}
