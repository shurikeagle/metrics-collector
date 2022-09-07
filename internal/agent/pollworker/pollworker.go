package pollworker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/agent/metric"
)

// Poller collects any metrics
type Poller interface {
	Poll() metric.Metrics
}

type pollWorker struct {
	running      bool
	pollInterval time.Duration
	currentStats metric.Metrics
	poller       Poller
	pollCounter  int64
	mx           sync.RWMutex
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
		currentStats: metric.Metrics{
			Gauges:   make(map[string]float64),
			Counters: make(map[string]int64),
		},
	}, nil
}

// Run starts pollWorker
func (w *pollWorker) Run(ctx context.Context) error {
	if w.running {
		return errors.New("poll worker is already running")
	}
	w.running = true

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.currentStats = w.poller.Poll()
			w.pollCounter++
		case <-ctx.Done():
			log.Println(ctx.Err(), ", stopping poll worker")
			w.running = false
			return nil
		}
	}
}

// Stats returns results of the last pollWorker's metrics poll
func (w *pollWorker) Stats() metric.Metrics {
	w.mx.RLock()
	defer w.mx.RUnlock()

	w.currentStats.Counters["PollCount"] = w.pollCounter
	w.pollCounter = 0

	return w.currentStats
}
