package metrichandler

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

var ErrUnexpectedMetricType = errors.New("unexpected metric type")
var ErrMetricNotFound = errors.New("metric not found")

// handler is a route handler for working with metrics through repository
type handler struct {
	*chi.Mux
	storage storage.MetricRepository
}

// New creates instance of metric handler
func New(s storage.MetricRepository) *handler {
	r := chi.NewRouter()
	h := &handler{
		Mux:     r,
		storage: s,
	}

	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)

	h.Get("/", h.getAllHandler())
	h.Get("/value/{metricType}/{metricName}", h.getValueHandler())

	h.Post("/update/{metricType}/{metricName}/{metricValue}", h.updateHandler())

	return h
}
