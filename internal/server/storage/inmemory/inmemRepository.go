package inmemory

import (
	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

var _ storage.MetricRepository = (*InmemMetricRepository)(nil)

// TODO: Chekc if need to initialize maps
type InmemMetricRepository struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (r *InmemMetricRepository) AddOrUpdateGauge(g metric.Gauge) {
	r.gauges[g.Name] = g.Value
}

func (r *InmemMetricRepository) AddOrUpdateCounter(c metric.Counter) {
	r.counters[c.Name] = c.Value
}

func (r *InmemMetricRepository) GetCounter(name string) (c metric.Counter, ok bool) {
	c.Name = name
	c.Value, ok = r.counters[c.Name]

	return
}
