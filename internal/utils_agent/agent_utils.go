package agent

import (
	"log"
	"net/http"
	rt "runtime"
	"strconv"
	"time"
)

type counter int

type config struct {
	pollInterval   time.Duration
	reportInterval time.Duration
}

// This function makes a copy of the map so that the metrics can be updated.
func copyMap(metricList *map[string]*uint64) map[string]uint64 {
	oldMetricList := make(map[string]uint64)
	for k, v := range *metricList {
		oldMetricList[k] = *v
	}
	return oldMetricList
}

// updating and counting updated metrics
func updateAndCount(m *rt.MemStats, metricList *map[string]*uint64, pollCounter *counter) {

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

// function to create a metrics list and Memory Stats obj
func createMetricList(pollCounter *counter) (metricList map[string]*uint64) {
	var m rt.MemStats
	rt.ReadMemStats(&m)

	GCCPUFraction := uint64(m.GCCPUFraction)
	metricList = map[string]*uint64{
		"Alloc":         &m.Alloc,
		"BuckHashSys":   &m.BuckHashSys,
		"Frees":         &m.Frees,
		"GCCPUFraction": &GCCPUFraction,
		"GCSys":         &m.GCSys,
		"HeapAlloc":     &m.HeapAlloc,
		"HeapIdle":      &m.HeapIdle,
	}

	updateAndCount(&m, &metricList, pollCounter)
	return
}

// sending collected metrics to 127.0.0.1/updates on port 8080
func sendData(metricsList *map[string]*uint64) {
	for k, v := range *metricsList {

		url := "http://127.0.0.1:8080/update/" + "type" + "/" + k + "/" + strconv.Itoa(int(*v)) + "/"
		_, err := http.Post(url, "text/plain", nil)
		if err != nil {
			log.Fatal(err)
		}

	}

}

// system monitor
func MyMonitor() {
	//metric counter
	var pollCount counter = 0
	//map for metric
	var metricsList map[string]*uint64

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
