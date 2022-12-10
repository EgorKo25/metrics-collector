package handlers

import (
	"fmt"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ShowThisMetricValue(m *storage.MetricStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := (*m).TakeThisStat(chi.URLParam(r, "name"), chi.URLParam(r, "type"))
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(fmt.Sprintf("%v\n", res)))
		if err != nil {
			log.Printf("Cannot write reqeust: %s", err)
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}
func ShowAllMetricFromStorage(m *storage.MetricStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for k, v := range (*m).MetricsGauge {
			tmp := "> " + k + ":  " + fmt.Sprintf("%f", v) + "\n"
			_, err := w.Write([]byte(tmp))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatalln(err)
			}
		}
		for k, v := range (*m).MetricsCounter {
			tmp := "> " + k + ":  " + fmt.Sprintf("%d", v) + "\n"
			_, err := w.Write([]byte(tmp))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatalln(err)
			}
		}
	}
}
func AddMetricToStorage(m *storage.MetricStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if mType := chi.URLParam(r, "type"); mType == "gauge" {
			value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("Somethings went wrong: %s", err)
				return
			}
			(*m).GetStats(chi.URLParam(r, "name"), any(storage.Gauge(value)), mType)
			w.WriteHeader(http.StatusOK)
		}
		if mType := chi.URLParam(r, "type"); mType == "counter" {
			value, err := strconv.Atoi(chi.URLParam(r, "value"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("Somethings went wrong: %s", err)
				return
			}

			(*m).GetStats(chi.URLParam(r, "name"), any(storage.Counter(value)), mType)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
	}
}

func GetMetricList(MetricList *map[string]storage.Gauge, CounterList *map[string]storage.Counter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		data := strings.Split(url, "/")

		if data[2] == "gauge" {
			if len(data) > 4 {
				value, err := strconv.ParseFloat(data[4], 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
				(*MetricList)[data[3]] = storage.Gauge(value)
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
				(*CounterList)[data[3]] = storage.Counter(value)
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
