package runtimePoller

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/metric"
	"github.com/shurikeagle/metrics-collector/internal/pollworker"
)

var _ pollworker.Poller = (*Poller)(nil)

// Poller is an object for collecting runtime metrics
type Poller struct {
	pollCounter int64
}

// Poll collects runtime metrics
func (p *Poller) Poll() metric.Metrics {
	p.pollCounter++

	m := metric.Metrics{
		Gauges:   make(map[string]float64, 28),
		Counters: make(map[string]int64, 1),
	}

	addNonRuntimeMetrics(&m, p.pollCounter)
	addRuntimeMetrics(&m)

	return m
}

func addNonRuntimeMetrics(m *metric.Metrics, poolCount int64) {
	rnd := rand.NewSource(time.Now().UnixNano())
	m.Gauges["RandomValue"] = float64(rnd.Int63())

	m.Counters["PollCount"] = poolCount
}

func addRuntimeMetrics(m *metric.Metrics) {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	m.Gauges["Alloc"] = float64(memstats.Alloc)
	m.Gauges["BuckHashSys"] = float64(memstats.BuckHashSys)
	m.Gauges["Frees"] = float64(memstats.Frees)
	m.Gauges["GCCPUFraction"] = float64(memstats.GCCPUFraction)
	m.Gauges["GCSys"] = float64(memstats.GCSys)
	m.Gauges["HeapAlloc"] = float64(memstats.HeapAlloc)
	m.Gauges["HeapIdle"] = float64(memstats.HeapIdle)
	m.Gauges["HeapInuse"] = float64(memstats.HeapInuse)
	m.Gauges["HeapObjects"] = float64(memstats.HeapObjects)
	m.Gauges["HeapReleased"] = float64(memstats.HeapReleased)
	m.Gauges["HeapSys"] = float64(memstats.HeapSys)
	m.Gauges["LastGC"] = float64(memstats.LastGC)
	m.Gauges["Lookups"] = float64(memstats.Lookups)
	m.Gauges["MCacheInuse"] = float64(memstats.MCacheInuse)
	m.Gauges["MCacheSys"] = float64(memstats.MCacheSys)
	m.Gauges["MSpanInuse"] = float64(memstats.MSpanInuse)
	m.Gauges["MSpanSys"] = float64(memstats.MSpanSys)
	m.Gauges["Mallocs"] = float64(memstats.Mallocs)
	m.Gauges["NextGC"] = float64(memstats.NextGC)
	m.Gauges["NumForcedGC"] = float64(memstats.NumForcedGC)
	m.Gauges["NumGC"] = float64(memstats.NumGC)
	m.Gauges["OtherSys"] = float64(memstats.OtherSys)
	m.Gauges["PauseTotalNs"] = float64(memstats.PauseTotalNs)
	m.Gauges["StackInuse"] = float64(memstats.StackInuse)
	m.Gauges["StackSys"] = float64(memstats.StackSys)
	m.Gauges["Sys"] = float64(memstats.Sys)
	m.Gauges["TotalAlloc"] = float64(memstats.TotalAlloc)
}
