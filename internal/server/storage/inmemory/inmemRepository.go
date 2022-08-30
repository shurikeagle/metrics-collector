package inmemory

import (
	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

var _ storage.MetricRepository = (*InmemMetricRepository)(nil)

type InmemMetricRepository struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (r *InmemMetricRepository) AddOrUpdateGauge(g metric.Gauge) {
	r.gauges[g.Name] = g.Value
}

func (r *InmemMetricRepository) AddOrUpdateCounter(c metric.Counter) {
	if curentVal, ok := r.counters[c.Name]; !ok {
		r.counters[c.Name] = curentVal
	} else {
		r.counters[c.Name] += curentVal
	}
}
