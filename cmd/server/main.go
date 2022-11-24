package main

import (
	"handlers/handlers"
	"log"
	"net/http"
)

func main() {
	MetricList := make(map[string]handlers.Guage)
	CounterList := make(map[string]handlers.Counter)

	http.HandleFunc("/update/", handlers.GetMetricList(&MetricList, &CounterList))

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	log.Fatal(server.ListenAndServe())
}
