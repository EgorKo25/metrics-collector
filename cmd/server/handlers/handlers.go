package server

import (
	"net/http"
	"strconv"
)

type Guage float64
type Counter uint64

func GetMetricList(MetricList *map[string]Guage, CounterList *map[string]Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataType := r.URL.Query().Get("type"); dataType == "Guage" {
			value, err := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			(*MetricList)[r.URL.Query().Get("name")] = Guage(value)

			w.WriteHeader(http.StatusOK)
			return

		}
		if dataType := r.URL.Query().Get("type"); dataType == "Counter" {
			value, err := strconv.Atoi(r.URL.Query().Get("value"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			(*CounterList)[r.URL.Query().Get("name")] = Counter(value)

			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
