package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Gauge float64
type Counter uint64

func GetMetricList(MetricList *map[string]Gauge, CounterList *map[string]Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		data := strings.Split(url, "/")

		if data[2] == "gauge" {
			if len(data) > 4 {
				value, err := strconv.ParseFloat(data[4], 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				(*MetricList)[data[3]] = Gauge(value)
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
				(*CounterList)[data[3]] = Counter(value)
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
