package routers

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(m *storage.MetricStorage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.ShowAllMetricFromStorage(m))
		r.Get("/value/{type}/{name}", handlers.ShowThisMetricValue(m))
		r.Post("/update/{type}/{name}/{value}", handlers.AddMetricToStorage(m))

	})

	return r
}
