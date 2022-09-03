package metrichandler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

// Handler is a metric handler for workin with metrics
type handler struct {
	*chi.Mux
	storage storage.MetricRepository
}

// New creates instance of metric handler
func New(s storage.MetricRepository) *handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := &handler{
		Mux:     chi.NewMux(),
		storage: s,
	}

	h.Post("/update/{metricType}/{metricName}/{metricValue}", h.updateHandler())

	return h
}
