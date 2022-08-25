package collectWorker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/metric"
)

const defaultPollInterval = 2 * time.Second

type Collector interface {
	Collect() metric.Metrics
}

type collectWorker struct {
	collectInterval time.Duration
	currentStats    metric.Metrics
	collector       Collector
}

func New(pollInterval time.Duration, collector Collector) (*collectWorker, error) {
	if pollInterval == 0 {
		pollInterval = defaultPollInterval
	}

	if collector == nil {
		return nil, errors.New("collector param is empty")
	}

	return &collectWorker{
		collectInterval: pollInterval,
		collector:       collector,
	}, nil
}

func (w *collectWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.collectInterval)

	for {
		select {
		case <-ticker.C:
			w.currentStats = w.collector.Collect()

			// TODO: Debug, to remove
			log.Println("============================")
			for k, v := range w.currentStats.Gauges {
				log.Println(k, ":", v)
			}
			for k, v := range w.currentStats.Counters {
				log.Println(k, ":", v)
			}
			log.Println("============================")
		case <-ctx.Done():
			log.Println(ctx.Err(), ", stopping collect worker")
			return
		}
	}
}

func (w *collectWorker) Stats() metric.Metrics {
	return w.currentStats
}
