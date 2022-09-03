package storage

import "github.com/shurikeagle/metrics-collector/internal/server/metric"

type MetricRepository interface {
	GetCounter(name string) (c metric.Counter, ok bool)
	GetGauge(name string) (c metric.Gauge, ok bool)
	AddOrUpdateGauge(g metric.Gauge)
	AddOrUpdateCounter(c metric.Counter)
}
