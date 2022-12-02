package handlers

import (
	"DevOps-Track-Yandex/internal/StorageSupport "
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetMetricList(MetricList *map[string]StorageSupport.Gauge, CounterList *map[string]StorageSupport.Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		data := strings.Split(url, "/")

		if data[2] == "gauge" {
			if len(data) > 4 {
				value, err := strconv.ParseFloat(data[4], 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				(*MetricList)[data[3]] = StorageSupport.Gauge(value)
				log.Println(MetricList)

				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusNotFound)
			return

		}
		if data[2] == "counter" {
			if len(data) > 4 {
				value, err := strconv.Atoi(data[4])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				(*CounterList)[data[3]] = StorageSupport.Counter(value)
				log.Println(CounterList)

				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return

		}
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}
func TakeDaefaultPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
