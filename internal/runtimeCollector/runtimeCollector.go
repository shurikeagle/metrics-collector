package runtime

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/collectWorker"
	"github.com/shurikeagle/metrics-collector/internal/metric"
)

var counter int64 = 0
var memstats *runtime.MemStats

var _ collectWorker.Collector = (*RuntimeCollector)(nil)

type RuntimeCollector struct{}

func (c RuntimeCollector) Collect() metric.Metrics {
	runtime.ReadMemStats(memstats)

	return createMetrics(*memstats)
}

func createMetrics(runtime.MemStats) metric.Metrics {
	m := metric.Metrics{
		Gauges:   make(map[string]float64, 0), // TODO: Presize
		Counters: make(map[string]int64, 0),   // TODO: Presize
	}

	// TODO: Move to collectWorker
	rnd := rand.NewSource(time.Now().UnixNano())
	m.Gauges["RandomValue"] = float64(rnd.Int63())

	counter++
	m.Counters["PollCount"] = counter

	return m
}
