package main

import (
	"DevOps-Track-Yandex/internal/ServerSupport/handlers"
	"log"
	"net/http"
)

func SimpleServerFunc() {

	http.HandleFunc("/update/", handlers.GetMetricList(&MetricList, &CounterList))

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	log.Fatal(server.ListenAndServe())
}
