package main

import (
	"log"
	"net/http"
	handlers "server/handlers"
)

func main() {
	var MetricList map[string]handlers.Guage
	var CounterList map[string]handlers.Counter

	http.HandleFunc("/update/", handlers.GetMetricList(&MetricList, &CounterList))

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	log.Fatal(server.ListenAndServe())
}
