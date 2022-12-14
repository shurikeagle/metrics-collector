package inmemory

import (
	"sync"

	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

var _ storage.MetricRepository = (*inmemMetricRepository)(nil)

type inmemMetricRepository struct {
	gauges   map[string]float64
	counters map[string]int64
	mx       sync.RWMutex
}

func New() *inmemMetricRepository {
	return &inmemMetricRepository{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (r *inmemMetricRepository) GetAll() ([]metric.Counter, []metric.Gauge) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	counters := make([]metric.Counter, 0, len(r.counters))
	gauges := make([]metric.Gauge, 0, len(r.gauges))

	for n, v := range r.counters {
		counters = append(counters, metric.Counter{
			Name:  n,
			Value: v,
		})
	}

	for n, v := range r.gauges {
		gauges = append(gauges, metric.Gauge{
			Name:  n,
			Value: v,
		})
	}

	return counters, gauges
}

func (r *inmemMetricRepository) GetCounter(name string) (c metric.Counter, ok bool) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	c.Name = name
	c.Value, ok = r.counters[c.Name]

	return
}

func (r *inmemMetricRepository) GetGauge(name string) (c metric.Gauge, ok bool) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	c.Name = name
	c.Value, ok = r.gauges[c.Name]

	return
}

func (r *inmemMetricRepository) AddOrUpdateGauge(g metric.Gauge) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.gauges[g.Name] = g.Value
}

func (r *inmemMetricRepository) AddOrUpdateCounter(c metric.Counter) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.counters[c.Name] = c.Value
}
