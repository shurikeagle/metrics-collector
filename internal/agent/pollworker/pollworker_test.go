package pollworker

import (
	"context"
	"testing"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/agent/metric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type moqPoller struct{}

var _ Poller = (*moqPoller)(nil)

func (p *moqPoller) Poll() metric.Metrics {
	return metric.Metrics{
		Gauges: map[string]float64{
			"gauge": 1.1,
		},
		Counters: map[string]int64{
			"counter": 1,
		},
	}
}

func Test_pollWorker_Run(t *testing.T) {
	worker, err := New(&moqPoller{}, 1*time.Second)
	require.NoError(t, err)

	ctx, cancelFunc := context.WithCancel(context.Background())
	go worker.Run(ctx)

	time.Sleep(1500 * time.Millisecond)
	cancelFunc()

	stats := worker.Stats()
	g, ok := stats.Gauges["gauge"]
	assert.Equal(t, true, ok)
	assert.Equal(t, float64(1.1), g)
	c, ok := stats.Counters["counter"]
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(1), c)
	pc, ok := stats.Counters["PollCount"]
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(1), pc)

	stats = worker.Stats()
	pc, ok = stats.Counters["PollCount"]
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(0), pc)
}
