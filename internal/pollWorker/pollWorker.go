package pollWorker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/metric"
)

// Poller collects any metrics
type Poller interface {
	Poll() metric.Metrics
}

type pollWorker struct {
	pollInterval time.Duration
	currentStats metric.Metrics
	poller       Poller
}

// New creates new instance of pollWorker.
// pollWorker is a job for collecting metrics with pollInterval
func New(poller Poller, pollInterval time.Duration) (*pollWorker, error) {
	if poller == nil {
		return nil, errors.New("collector param is empty")
	}

	if pollInterval == 0 {
		return nil, errors.New("pollInterval can't be 0")
	}

	return &pollWorker{
		pollInterval: pollInterval,
		poller:       poller,
	}, nil
}

// Run starts pollWorker
func (w *pollWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)

	for {
		select {
		case <-ticker.C:
			w.currentStats = w.poller.Poll()
		case <-ctx.Done():
			log.Println(ctx.Err(), ", stopping poll worker")
			return
		}
	}
}

// Stats returns results of the last pollWorker's metrics poll
func (w *pollWorker) Stats() metric.Metrics {
	return w.currentStats
}
