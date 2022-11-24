package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Guage float64
type Counter uint64

func GetMetricList(MetricList *map[string]Guage, CounterList *map[string]Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		data := strings.Split(url, "/")

		if data[2] == "guage" {
			if len(data) > 4 {
				value, err := strconv.ParseFloat(data[4], 64)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				(*MetricList)[data[3]] = Guage(value)
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
					w.WriteHeader(http.StatusInternalServerError)
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