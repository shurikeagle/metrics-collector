package metrichandler

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

const JSONcontentType = "application/json"

const JSONcontentTypeExpected = "expected 'application/json' content type"

var ErrInvalidRequestBody = errors.New("invalid request body")
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

	// Get
	h.Get("/", h.getAllHandler)
	h.Get("/value/{metricType}/{metricName}", h.getValueFromPathHandler)
	h.Post("/value", h.getValueFromBodyHandler)
	h.Post("/value/", h.getValueFromBodyHandler)

	// Update
	h.Post("/update", h.updateHandlerFromBody)
	h.Post("/update/", h.updateHandlerFromBody)
	h.Post("/update/{metricType}/{metricName}/{metricValue}", h.updateHandlerFromPath)

	return h
}
