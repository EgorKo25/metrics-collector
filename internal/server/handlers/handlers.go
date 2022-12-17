package handlers

import (
	"fmt"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	storage *storage.MetricStorage
}

// NewHandler handler type constructor
func NewHandler(storage *storage.MetricStorage) *Handler {
	return &Handler{
		storage: storage,
	}
}

// GetValueStat a handler that returns the value of a specific metric
func (h Handler) GetValueStat(w http.ResponseWriter, r *http.Request) {
	res := h.storage.StatStatus(chi.URLParam(r, "name"), chi.URLParam(r, "type"))
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

// GetAllStats returns the values of all metrics
func (h Handler) GetAllStats(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	for k, v := range h.storage.MetricsGauge {

		tmp := "> " + k + ":  " + fmt.Sprintf("%f", v) + "\n"

		_, err := w.Write([]byte(tmp))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalln(err)
		}
	}

	for k, v := range h.storage.MetricsCounter {

		tmp := "> " + k + ":  " + fmt.Sprintf("%d", v) + "\n"

		_, err := w.Write([]byte(tmp))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalln(err)
		}
	}
}

// SetMetricValue sets the value of the specified metric
func (h Handler) SetMetricValue(w http.ResponseWriter, r *http.Request) {

	if mType := chi.URLParam(r, "type"); mType == "gauge" {

		value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Somethings went wrong: %s", err)
			return
		}

		h.storage.SetGaugeStat(chi.URLParam(r, "name"), storage.Gauge(value), mType)
		w.WriteHeader(http.StatusOK)
	}
	if mType := chi.URLParam(r, "type"); mType == "counter" {

		value, err := strconv.Atoi(chi.URLParam(r, "value"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Somethings went wrong: %s", err)
			return
		}

		h.storage.SetCounterStat(chi.URLParam(r, "name"), storage.Counter(value), mType)
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusNotImplemented)
	return

}
