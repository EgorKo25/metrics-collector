package routers

import (
	"DevOps-Track-Yandex/internal/ServerSupport/handlers"
	"DevOps-Track-Yandex/internal/StorageSupport"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(m StorageSupport.MemStats) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.ShowAllMetricFromStorage(m))
		r.Get("/value/{type}/{name}", handlers.ShowThisMetricValue(m, r))
		r.Post("/update/{type}/{name}/{value}", handlers.AddMetricToStorage(m, r))
		r.Post("/update", handlers.GetMetricList(&m.MetricsGauge, &m.MetricsCounter))

	})

	return r
}
