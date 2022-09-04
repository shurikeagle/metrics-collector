package storage

import "github.com/shurikeagle/metrics-collector/internal/server/metric"

type MetricRepository interface {
	GetAll() ([]metric.Counter, []metric.Gauge)
	GetCounter(name string) (c metric.Counter, ok bool)
	GetGauge(name string) (c metric.Gauge, ok bool)
	AddOrUpdateCounter(c metric.Counter)
	AddOrUpdateGauge(g metric.Gauge)
}
