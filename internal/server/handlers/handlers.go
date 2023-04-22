// Package handlers описывает структуру хэндлеров и все хэндлеры веб приложения
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/database"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/encryption"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/hashing"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/middleware"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
)

// Handler структура определяющая структуру обработчиков
type Handler struct {
	storage    *storage.MetricStorage
	compressor *middleware.Compressor
	hasher     *hashing.Hash
	db         *database.DB
	encryptor  *encryption.Encryptor
}

// NewHandler конструктор структуры хэндлер
func NewHandler(storage *storage.MetricStorage, compressor *middleware.Compressor, hasher *hashing.Hash, db *database.DB, enc *encryption.Encryptor) *Handler {

	return &Handler{
		storage:    storage,
		compressor: compressor,
		hasher:     hasher,
		db:         db,
		encryptor:  enc,
	}
}

// PingDB проверяет соединение с базой данных
func (h *Handler) PingDB(w http.ResponseWriter, _ *http.Request) {

	ctx := context.Background()

	if err := h.db.DB.PingContext(ctx); err != nil {
		log.Println("database didn't open")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// GetValueStat возвращает значение конкретной метрики
func (h *Handler) GetValueStat(w http.ResponseWriter, r *http.Request) {
	res := h.storage.StatStatusM(chi.URLParam(r, "name"), chi.URLParam(r, "type"))
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

// GetJSONValue возвращает значение метрики в формате JSON
func (h *Handler) GetJSONValue(w http.ResponseWriter, r *http.Request) {

	var err error
	var metric storage.Metric

	b, _ := io.ReadAll(r.Body)

	b, err = h.encryptor.Decrypt(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(b, &metric)
	if err != nil {
		fmt.Printf("Unmarshal went wrong:  %s\n", err)
	}

	stat := h.storage.StatStatusM(metric.ID, metric.MType)

	switch metric.MType {

	case "gauge":
		if stat == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if stat != nil {

			tmp := stat.(storage.Gauge)

			metric.Value = &tmp
			metric.Delta = nil

		}
	case "counter":
		if stat == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if stat != nil {

			tmp := stat.(storage.Counter)

			metric.Delta = &tmp
			metric.Value = nil

		}
	}

	if metric.Hash, err = h.hasher.Run(&metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dataJSON, err := json.Marshal(metric)
	if err != nil {
		log.Println("failed to serialize!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept-Encoding") == "gzip" {

		dataJSON, err = h.compressor.Compress(dataJSON)
		if err != nil {
			log.Println("failed to compress")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Encoding", "gzip")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(dataJSON)
	return

}

// GetJSONUpdates обновляет пакет метрик принимая их из формата JSON
func (h *Handler) GetJSONUpdates(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	var err error
	var Metrics []storage.Metric

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read request body error!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err = h.encryptor.Decrypt(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Encoding") == "gzip" {

		b, err = h.compressor.Decompress(b)
		if err != nil {

			log.Println("field to decompress!")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = json.Unmarshal(b, &Metrics)
	if err != nil {

		log.Println("unmarshal went wrong!")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, metric := range Metrics {
		if metric.Value == nil && metric.Delta == nil {

			log.Println("no metric value")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if metric.Hash, err = h.hasher.Run(&metric); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.storage.SetStat(&metric)
		if err = h.addMetric(ctx, &metric); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
}

// addMetric добавялет метрику в буффер, при заполнении буфера создает транзакцию в бд
func (h *Handler) addMetric(ctx context.Context, m *storage.Metric) error {

	h.db.Buffer = append(h.db.Buffer, *m)

	if cap(h.db.Buffer) == len(h.db.Buffer) {
		err := h.db.FlushWithContext(ctx)
		if err != nil {
			return errors.New("cannot add records to the database")
		}
	}
	return nil
}

// SetJSONValue устанавливает значение метрики, принимает формат JSON
func (h *Handler) SetJSONValue(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	var err error
	var metric storage.Metric

	b, _ := io.ReadAll(r.Body)

	b, err = h.encryptor.Decrypt(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Encoding") == "gzip" {
		b, _ = h.compressor.Decompress(b)
	}

	err = json.Unmarshal(b, &metric)
	if err != nil {
		fmt.Printf("Unmarshal went wrong:  %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if metric.Value == nil && metric.Delta == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric.Hash, err = h.hasher.Run(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":

		if metric.Value == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if metric.Value != nil {
			h.storage.SetStat(&metric)
		}

		if stat := h.storage.StatStatusM(metric.ID, metric.MType); stat != nil {
			tmp := stat.(storage.Gauge)
			metric.Value = &tmp
		}

	case "counter":

		if metric.Delta == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if metric.Delta != nil {
			h.storage.SetStat(&metric)
		}

		if stat := h.storage.StatStatusM(metric.ID, metric.MType); stat != nil && stat.(storage.Counter) != 0 {
			tmp := stat.(storage.Counter)
			metric.Delta = &tmp
		}

	}

	if h.db != nil {
		if err = h.db.Run(ctx, &metric); err != nil {
			log.Println("Error db send ", err)
		}
	}

	dataJSON, err := json.Marshal(metric)
	if err != nil {
		log.Println("Failed to serialize")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept-Encoding") == "gzip" {

		dataJSON, err = h.compressor.Compress(dataJSON)
		if err != nil {
			log.Println("Failed to middleware")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Encoding", "gzip")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(dataJSON)
	return

}

// GetAllStats возвращает значения всех метрик
func (h *Handler) GetAllStats(w http.ResponseWriter, r *http.Request) {

	var res string
	var err error

	for k, v := range h.storage.Metrics {

		if v.MType == "gauge" {
			res += "> " + k + ":  " + fmt.Sprintf("%f", *v.Value) + "\n"
		}
		if v.MType == "counter" {
			res += "> " + k + ":  " + fmt.Sprintf("%d", *v.Delta) + "\n"
		}

	}

	tmp := []byte(res)

	if r.Header.Get("Accept-Encoding") == "gzip" {
		tmp, err = h.compressor.Compress(tmp)
		if err != nil {
			log.Println("Failed to compress")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Encoding", "gzip")
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(tmp)
}

// SetMetricValue устанавливает значение метрики переданное в url
func (h *Handler) SetMetricValue(w http.ResponseWriter, r *http.Request) {

	var metric storage.Metric

	if mType := chi.URLParam(r, "type"); mType == "gauge" {

		tmp, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Somethings went wrong: %s", err)
			return
		}
		value := storage.Gauge(tmp)

		metric.ID = chi.URLParam(r, "name")
		metric.MType = mType
		metric.Value = &value
		metric.Delta = nil

		h.storage.SetStat(&metric)
		w.WriteHeader(http.StatusOK)
	}

	if mType := chi.URLParam(r, "type"); mType == "counter" {

		tmp, err := strconv.Atoi(chi.URLParam(r, "value"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Somethings went wrong: %s", err)
			return
		}
		value := storage.Counter(tmp)

		metric.ID = chi.URLParam(r, "name")
		metric.MType = mType
		metric.Value = nil
		metric.Delta = &value

		h.storage.SetStat(&metric)
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusNotImplemented)
	return

}
