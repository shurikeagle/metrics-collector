package collectWorker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/metric"
)

const defaultPollInterval = 2 * time.Second // TODO: move to main

// Collector collects any metrics
type Collector interface {
	Collect() metric.Metrics
}

type collectWorker struct {
	collectInterval time.Duration
	currentStats    metric.Metrics
	collector       Collector
}

// New creates new instance of collectWorker.
// collectWorker is a job for collecting metrics with poolInterval
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

// Run starts collectWorker
func (w *collectWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.collectInterval)

	for {
		select {
		case <-ticker.C:
			w.currentStats = w.collector.Collect()
		case <-ctx.Done():
			log.Println(ctx.Err(), ", stopping collect worker")
			return
		}
	}
}

// Stats returns results of the last collectWorker's metrics pool
func (w *collectWorker) Stats() metric.Metrics {
	return w.currentStats
}
