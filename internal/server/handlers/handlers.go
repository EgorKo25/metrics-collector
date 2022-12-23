package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/serializer"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	storage    *storage.MetricStorage
	serializer *serializer.Serialize
}

// NewHandler handler type constructor
func NewHandler(storage *storage.MetricStorage, srl *serializer.Serialize) *Handler {
	return &Handler{
		storage:    storage,
		serializer: srl,
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

// GetJSONValue go dock
func (h Handler) GetJSONValue(w http.ResponseWriter, r *http.Request) {

	h.serializer.Clean()
	w.Header().Add("Content-Type", "application/json")
	b, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(b, h.serializer); err != nil {
		fmt.Printf("Unmarshal went wrong:  %s\n", err)
	}

	stat := h.storage.StatStatus(h.serializer.ID, h.serializer.MType)

	switch h.serializer.MType {
	case "gauge":
		if stat == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if stat != nil {
			tmp := stat.(storage.Gauge)
			h.serializer.Value = &tmp
			h.serializer.Delta = nil
			log.Printf("%f, %s, %s", *h.serializer.Value, h.serializer.ID, h.serializer.MType)
		}
	case "counter":
		if stat == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if stat != nil {
			tmp := stat.(storage.Counter)
			h.serializer.Delta = &tmp
			h.serializer.Value = nil
			log.Printf(" In Block Counter: %d, %s, %s", *h.serializer.Delta, h.serializer.ID, h.serializer.MType)
		}
	}

	if dataJSON, err := h.serializer.Run(); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write(dataJSON)
	}

}

// SetJSONValue go dock
func (h Handler) SetJSONValue(w http.ResponseWriter, r *http.Request) {
	h.serializer.Clean()
	w.Header().Add("content-Type", "application/json")
	b, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(b, h.serializer); err != nil {
		fmt.Printf("Unmarshal went wrong:  %s\n", err)
	}

	if h.serializer.Value == nil && h.serializer.Delta == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch h.serializer.MType {
	case "gauge":
		if h.serializer.Delta == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if h.serializer.Value != nil {
			log.Printf("If not nil %f, %s, %s", *h.serializer.Value, h.serializer.ID, h.serializer.MType)
			h.storage.SetGaugeStat(h.serializer.ID, *h.serializer.Value, h.serializer.MType)

		}
		if stat := h.storage.StatStatus(h.serializer.ID, h.serializer.MType); stat != nil {
			log.Printf("In Block Guage: %f, %s, %s", *h.serializer.Value, h.serializer.ID, h.serializer.MType)
			tmp := stat.(storage.Gauge)
			h.serializer.Value = &tmp
		}
	case "counter":
		if h.serializer.Delta == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if h.serializer.Delta != nil {
			h.storage.SetCounterStat(h.serializer.ID, *h.serializer.Delta, h.serializer.MType)
		}
		if stat := h.storage.StatStatus(h.serializer.ID, h.serializer.MType); stat != nil && stat.(storage.Counter) != 0 {
			log.Printf("In Block Counter: %d, %s, %s", *h.serializer.Delta, h.serializer.ID, h.serializer.MType)
			tmp := stat.(storage.Counter)
			h.serializer.Delta = &tmp
		}

	}

	if dataJSON, err := h.serializer.Run(); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-Type", "application/json")
		_, _ = w.Write(dataJSON)
	}

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
