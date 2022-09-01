package metrichandler

import (
	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

// Handler is a metric handler for workin with metrics
type Handler struct {
	storage storage.MetricRepository
}

// New creates instance of metric handler
func New(s storage.MetricRepository) *Handler {
	return &Handler{
		storage: s,
	}
}

// Update updates gauge metric in storage
func (h *Handler) UpdateGauge(g metric.Gauge) {
	h.storage.AddOrUpdateGauge(g)
}

// Update updates counter metric in storage
func (h *Handler) UpdateCounter(c metric.Counter) {
	existingCounter, _ := h.storage.GetCounter(c.Name)
	existingCounter.Value += c.Value

	h.storage.AddOrUpdateCounter(existingCounter)
}
