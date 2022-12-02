package handlers

import (
	"DevOps-Track-Yandex/internal/StorageSupport"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ShowThisMetricValue(m StorageSupport.MemStats, r chi.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := m.TakeThisStat(chi.URLParam(r, "name"), chi.URLParam(r, "type"))
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v\n", res)))
		w.WriteHeader(http.StatusOK)
		return

	}
}
func ShowAllMetricFromStorage(m StorageSupport.MemStats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for k, v := range m.MetricsGauge {
			tmp := "> " + k + ":  " + fmt.Sprintf("%f", v) + "\n"
			w.Write([]byte(tmp))
		}
		for k, v := range m.MetricsCounter {
			tmp := "> " + k + ":  " + fmt.Sprintf("%i", v) + "\n"
			w.Write([]byte(tmp))
		}
	}
}
func AddMetricToStorage(m StorageSupport.MemStats, r chi.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if mType := chi.URLParam(r, "type"); mType == "gauge" {
			value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println("Somethings went wrong: %s", err)
				return
			}
			m.GetStats(chi.URLParam(r, "name"), any(StorageSupport.Gauge(value)), mType)
			w.WriteHeader(http.StatusOK)
		}
		if mType := chi.URLParam(r, "type"); mType == "counter" {
			value, err := strconv.Atoi(chi.URLParam(r, "value"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println("Somethings went wrong: %s", err)
				return
			}

			m.GetStats(chi.URLParam(r, "name"), any(StorageSupport.Counter(value)), mType)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

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
		w.Write([]byte("SOSI"))
	}
}
