package simple_server

import (
	"DevOps-Track-Yandex/internal/ServerSupport/handlers"
	"DevOps-Track-Yandex/internal/StorageSupport"
	"log"
	"net/http"
)

func SimpleServerFunc() {

	MetricList := make(map[string]StorageSupport.Gauge)
	CounterList := make(map[string]StorageSupport.Counter)

	http.HandleFunc("/update/", handlers.GetMetricList(&MetricList, &CounterList))

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	log.Fatal(server.ListenAndServe())
}
