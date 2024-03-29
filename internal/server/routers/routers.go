// Package routers описывает роутер
package routers

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	mid "github.com/EgorKo25/DevOps-Track-Yandex/internal/server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter создет роутер определяющий маршруты к обработчикам
func NewRouter(handler *handlers.Handler, middle *mid.Middle) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middle.IpChecker)

	r.Mount("/debug", middleware.Profiler())

	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetAllStats)
		r.Get("/ping", handler.PingDB)
		r.Get("/value/{type}/{name}", handler.GetValueStat)
		r.Post("/update/{type}/{name}/{value}", handler.SetMetricValue)

		r.Post("/updates/", handler.GetJSONUpdates)
		r.Post("/update/", handler.SetJSONValue)
		r.Post("/value/", handler.GetJSONValue)

	})

	return r
}
