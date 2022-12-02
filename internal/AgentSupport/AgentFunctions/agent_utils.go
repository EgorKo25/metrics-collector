package AgentFunctions

import (
	"DevOps-Track-Yandex/internal/StorageSupport"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	rt "runtime"
	"time"
)

type config struct {
	pollInterval   time.Duration
	reportInterval time.Duration
}

// This function makes a copy of the map so that the metrics can be updated.
func copyMap(metricList *map[string]*StorageSupport.Gauge) map[string]StorageSupport.Gauge {
	oldMetricList := make(map[string]StorageSupport.Gauge)
	for k, v := range *metricList {
		oldMetricList[k] = *v
	}
	return oldMetricList
}

// updating and counting updated metrics
func updateAndCount(m *rt.MemStats, metricList *map[string]*StorageSupport.Gauge, pollCounter *StorageSupport.Counter) {

	(*metricList)["RandomValue"] = metricPtrCreator(StorageSupport.Gauge(rand.Float64()))
	oldMetricList := copyMap(metricList)
	//Updating
	rt.ReadMemStats(m)
	//Counting
	for k := range *metricList {
		if oldMetricList[k] != *(*metricList)[k] {
			*pollCounter++
		}
	}
}
func metricPtrCreator(val StorageSupport.Gauge) *StorageSupport.Gauge {
	return &val
}

// function to create a metrics list and Memory Stats obj
func createMetricList(pollCounter *StorageSupport.Counter) (metricList map[string]*StorageSupport.Gauge) {
	var m rt.MemStats
	rt.ReadMemStats(&m)

	metricList = make(map[string]*StorageSupport.Gauge)
	metricList["Alloc"] = metricPtrCreator(StorageSupport.Gauge(m.Alloc))
	metricList["BuckHashSys"] = metricPtrCreator(StorageSupport.Gauge(m.BuckHashSys))
	metricList["Frees"] = metricPtrCreator(StorageSupport.Gauge(m.Frees))
	metricList["GCCPUFraction"] = metricPtrCreator(StorageSupport.Gauge(m.GCCPUFraction))
	metricList["GCSys"] = metricPtrCreator(StorageSupport.Gauge(m.GCSys))
	metricList["HeapAlloc"] = metricPtrCreator(StorageSupport.Gauge(m.HeapAlloc))
	metricList["HeapIdle"] = metricPtrCreator(StorageSupport.Gauge(m.HeapIdle))
	metricList["HeapInuse"] = metricPtrCreator(StorageSupport.Gauge(m.HeapInuse))
	metricList["HeapObjects"] = metricPtrCreator(StorageSupport.Gauge(m.HeapObjects))
	metricList["HeapReleased"] = metricPtrCreator(StorageSupport.Gauge(m.HeapReleased))
	metricList["HeapSys"] = metricPtrCreator(StorageSupport.Gauge(m.HeapSys))
	metricList["LastGC"] = metricPtrCreator(StorageSupport.Gauge(m.LastGC))
	metricList["Lookups"] = metricPtrCreator(StorageSupport.Gauge(m.Lookups))
	metricList["MCacheInuse"] = metricPtrCreator(StorageSupport.Gauge(m.MCacheInuse))
	metricList["MCacheSys"] = metricPtrCreator(StorageSupport.Gauge(m.MCacheSys))
	metricList["MSpanInuse"] = metricPtrCreator(StorageSupport.Gauge(m.MSpanInuse))
	metricList["MSpanSys"] = metricPtrCreator(StorageSupport.Gauge(m.MSpanSys))
	metricList["Mallocs"] = metricPtrCreator(StorageSupport.Gauge(m.Mallocs))
	metricList["NextGC"] = metricPtrCreator(StorageSupport.Gauge(m.NextGC))
	metricList["NumForcedGC"] = metricPtrCreator(StorageSupport.Gauge(m.NumForcedGC))
	metricList["NumGC"] = metricPtrCreator(StorageSupport.Gauge(m.NumGC))
	metricList["OtherSys"] = metricPtrCreator(StorageSupport.Gauge(m.OtherSys))
	metricList["PauseTotalNs"] = metricPtrCreator(StorageSupport.Gauge(m.PauseTotalNs))
	metricList["StackInuse"] = metricPtrCreator(StorageSupport.Gauge(m.StackInuse))
	metricList["StackSys"] = metricPtrCreator(StorageSupport.Gauge(m.StackSys))
	metricList["Sys"] = metricPtrCreator(StorageSupport.Gauge(m.Sys))
	metricList["TotalAlloc"] = metricPtrCreator(StorageSupport.Gauge(m.TotalAlloc))
	metricList["RandomValue"] = metricPtrCreator(StorageSupport.Gauge(rand.Float64()))

	updateAndCount(&m, &metricList, pollCounter)
	return
}

// sending collected metrics to 127.0.0.1/updates on port 8080
func sendData(metricsList *map[string]*StorageSupport.Gauge) {
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

// system monitor
func MyMonitor() {
	//metric counter
	var pollCount StorageSupport.Counter = 0
	//map for metric
	var metricsList map[string]*StorageSupport.Gauge

	//interval settings
	var durConf config
	durConf.pollInterval = 2
	durConf.reportInterval = 10

	tickerInterval := time.NewTicker(durConf.pollInterval * time.Second)
	tickerReport := time.NewTicker(durConf.reportInterval * time.Second)

	for {
		select {
		case <-tickerInterval.C:
			metricsList = createMetricList(&pollCount)

		case <-tickerReport.C:
			sendData(&metricsList)
		}
	}

}
