package inmemory

import (
	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

var _ storage.MetricRepository = (*inmemMetricRepository)(nil)

type inmemMetricRepository struct {
	gauges   map[string]float64
	counters map[string]int64
}

func New() *inmemMetricRepository {
	return &inmemMetricRepository{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (r *inmemMetricRepository) AddOrUpdateGauge(g metric.Gauge) {
	r.gauges[g.Name] = g.Value
}

func (r *inmemMetricRepository) AddOrUpdateCounter(c metric.Counter) {
	r.counters[c.Name] = c.Value
}

func (r *inmemMetricRepository) GetCounter(name string) (c metric.Counter, ok bool) {
	c.Name = name
	c.Value, ok = r.counters[c.Name]

	return
}
