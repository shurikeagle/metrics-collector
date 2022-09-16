package runtimepoller

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/agent/metric"
	"github.com/shurikeagle/metrics-collector/internal/agent/pollworker"
)

var _ pollworker.Poller = (*Poller)(nil)

// Poller is an object for collecting runtime metrics
type Poller struct {
}

// Poll collects runtime metrics
func (p *Poller) Poll(m *metric.Metrics) {
	addNonRuntimeMetrics(m)
	addRuntimeMetrics(m)
}

func addNonRuntimeMetrics(m *metric.Metrics) {
	rnd := rand.NewSource(time.Now().UnixNano())
	m.SetGauge("RandomValue", float64(rnd.Int63()))
}

func addRuntimeMetrics(m *metric.Metrics) {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	m.SetGauge("Alloc", float64(memstats.Alloc))
	m.SetGauge("BuckHashSys", float64(memstats.BuckHashSys))
	m.SetGauge("Frees", float64(memstats.Frees))
	m.SetGauge("GCCPUFraction", float64(memstats.GCCPUFraction))
	m.SetGauge("GCSys", float64(memstats.GCSys))
	m.SetGauge("HeapAlloc", float64(memstats.HeapAlloc))
	m.SetGauge("HeapIdle", float64(memstats.HeapIdle))
	m.SetGauge("HeapInuse", float64(memstats.HeapInuse))
	m.SetGauge("HeapObjects", float64(memstats.HeapObjects))
	m.SetGauge("HeapReleased", float64(memstats.HeapReleased))
	m.SetGauge("HeapSys", float64(memstats.HeapSys))
	m.SetGauge("LastGC", float64(memstats.LastGC))
	m.SetGauge("Lookups", float64(memstats.Lookups))
	m.SetGauge("MCacheInuse", float64(memstats.MCacheInuse))
	m.SetGauge("MCacheSys", float64(memstats.MCacheSys))
	m.SetGauge("MSpanInuse", float64(memstats.MSpanInuse))
	m.SetGauge("MSpanSys", float64(memstats.MSpanSys))
	m.SetGauge("Mallocs", float64(memstats.Mallocs))
	m.SetGauge("NextGC", float64(memstats.NextGC))
	m.SetGauge("NumForcedGC", float64(memstats.NumForcedGC))
	m.SetGauge("NumGC", float64(memstats.NumGC))
	m.SetGauge("OtherSys", float64(memstats.OtherSys))
	m.SetGauge("PauseTotalNs", float64(memstats.PauseTotalNs))
	m.SetGauge("StackInuse", float64(memstats.StackInuse))
	m.SetGauge("StackSys", float64(memstats.StackSys))
	m.SetGauge("Sys", float64(memstats.Sys))
	m.SetGauge("TotalAlloc", float64(memstats.TotalAlloc))
}
