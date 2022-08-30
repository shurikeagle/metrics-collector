package storage

import "github.com/shurikeagle/metrics-collector/internal/server/metric"

type MetricRepository interface {
	AddOrUpdateGauge(g metric.Gauge)
	AddOrUpdateCounter(c metric.Counter)
}
