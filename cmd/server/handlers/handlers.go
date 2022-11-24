package server

import (
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
			value, err := strconv.ParseFloat(data[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			(*MetricList)[data[3]] = Guage(value)

			w.WriteHeader(http.StatusOK)
			return

		}
		if data[2] == "counter" {
			value, err := strconv.Atoi(data[4])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			(*CounterList)[data[3]] = Counter(value)

			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
