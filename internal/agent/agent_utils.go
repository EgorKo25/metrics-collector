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

func copyMap(metricList *map[string]*storage.Gauge) map[string]storage.Gauge {
	oldMetricList := make(map[string]storage.Gauge)
	for k, v := range *metricList {
		oldMetricList[k] = *v
	}
	return oldMetricList
}

func updateAndCount(m *runtime.MemStats, metricList *map[string]*storage.Gauge, pollCounter *storage.Counter) {

	(*metricList)["RandomValue"] = metricPtrCreator(storage.Gauge(rand.Float64()))
	oldMetricList := copyMap(metricList)
	//Updating
	runtime.ReadMemStats(m)
	//Counting
	for k := range *metricList {
		if oldMetricList[k] != *(*metricList)[k] {
			*pollCounter++
		}
	}
}
func metricPtrCreator(val storage.Gauge) *storage.Gauge {
	return &val
}

func createMetricList(pollCounter *storage.Counter) (metricList map[string]*storage.Gauge) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metricList = make(map[string]*storage.Gauge)
	metricList["Alloc"] = metricPtrCreator(storage.Gauge(m.Alloc))
	metricList["BuckHashSys"] = metricPtrCreator(storage.Gauge(m.BuckHashSys))
	metricList["Frees"] = metricPtrCreator(storage.Gauge(m.Frees))
	metricList["GCCPUFraction"] = metricPtrCreator(storage.Gauge(m.GCCPUFraction))
	metricList["GCSys"] = metricPtrCreator(storage.Gauge(m.GCSys))
	metricList["HeapAlloc"] = metricPtrCreator(storage.Gauge(m.HeapAlloc))
	metricList["HeapIdle"] = metricPtrCreator(storage.Gauge(m.HeapIdle))
	metricList["HeapInuse"] = metricPtrCreator(storage.Gauge(m.HeapInuse))
	metricList["HeapObjects"] = metricPtrCreator(storage.Gauge(m.HeapObjects))
	metricList["HeapReleased"] = metricPtrCreator(storage.Gauge(m.HeapReleased))
	metricList["HeapSys"] = metricPtrCreator(storage.Gauge(m.HeapSys))
	metricList["LastGC"] = metricPtrCreator(storage.Gauge(m.LastGC))
	metricList["Lookups"] = metricPtrCreator(storage.Gauge(m.Lookups))
	metricList["MCacheInuse"] = metricPtrCreator(storage.Gauge(m.MCacheInuse))
	metricList["MCacheSys"] = metricPtrCreator(storage.Gauge(m.MCacheSys))
	metricList["MSpanInuse"] = metricPtrCreator(storage.Gauge(m.MSpanInuse))
	metricList["MSpanSys"] = metricPtrCreator(storage.Gauge(m.MSpanSys))
	metricList["Mallocs"] = metricPtrCreator(storage.Gauge(m.Mallocs))
	metricList["NextGC"] = metricPtrCreator(storage.Gauge(m.NextGC))
	metricList["NumForcedGC"] = metricPtrCreator(storage.Gauge(m.NumForcedGC))
	metricList["NumGC"] = metricPtrCreator(storage.Gauge(m.NumGC))
	metricList["OtherSys"] = metricPtrCreator(storage.Gauge(m.OtherSys))
	metricList["PauseTotalNs"] = metricPtrCreator(storage.Gauge(m.PauseTotalNs))
	metricList["StackInuse"] = metricPtrCreator(storage.Gauge(m.StackInuse))
	metricList["StackSys"] = metricPtrCreator(storage.Gauge(m.StackSys))
	metricList["Sys"] = metricPtrCreator(storage.Gauge(m.Sys))
	metricList["TotalAlloc"] = metricPtrCreator(storage.Gauge(m.TotalAlloc))
	metricList["RandomValue"] = metricPtrCreator(storage.Gauge(rand.Float64()))

	updateAndCount(&m, &metricList, pollCounter)
	return
}

func sendData(metricsList *map[string]*storage.Gauge) {
	for k, v := range *metricsList {
		mType := "gauge"
		URL, err := url.JoinPath("http://127.0.0.1:8080", "update", mType, k, fmt.Sprintf("%f", *v))
		if err != nil {
			log.Fatalln(err)
		}

		_, err = http.Post(URL, "text/plain", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func MyMonitor() {
	//metric counter
	var pollCount storage.Counter = 0
	//map for metric
	var metricsList map[string]*storage.Gauge
	// interval cfg
	var durCfg = config.NewAgentConfig(2, 5)

	tickerPoll := time.NewTicker(durCfg.PollInterval * time.Second)
	tickerReport := time.NewTicker(durCfg.ReportInterval * time.Second)

	for {
		select {
		case <-tickerPoll.C:
			metricsList = createMetricList(&pollCount)

		case <-tickerReport.C:
			sendData(&metricsList)
		}
	}

}
