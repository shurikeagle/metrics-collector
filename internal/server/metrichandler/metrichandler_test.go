package metrichandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type moqMetricRepository struct{}

var _ storage.MetricRepository = (*moqMetricRepository)(nil)

func (r *moqMetricRepository) GetCounter(name string) (c metric.Counter, ok bool) {
	return metric.Counter{}, true
}
func (r *moqMetricRepository) GetGauge(name string) (c metric.Gauge, ok bool) {
	return metric.Gauge{}, true
}
func (r *moqMetricRepository) AddOrUpdateGauge(g metric.Gauge)     {}
func (r *moqMetricRepository) AddOrUpdateCounter(c metric.Counter) {}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{
			name:         "update without other path",
			method:       http.MethodPost,
			path:         "/update/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "update metric without value",
			method:       http.MethodPost,
			path:         "/update/Gauge/Name/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "update unexpected metric",
			method:       http.MethodPost,
			path:         "/update/Unexpected/Name/42",
			expectedCode: http.StatusNotImplemented,
		},
		{
			name:         "invalid hhtp method",
			method:       http.MethodGet,
			path:         "/update/Counter/Name/42",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "update gauge with string value",
			method:       http.MethodPost,
			path:         "/update/Counter/Name/42.1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "update counter with float value",
			method:       http.MethodPost,
			path:         "/update/Counter/Name/42.1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "correct counter update",
			method:       http.MethodPost,
			path:         "/update/Counter/Name/42",
			expectedCode: http.StatusOK,
		},
		{
			name:         "correct gauge update",
			method:       http.MethodPost,
			path:         "/update/Gauge/Name/42.1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "correct gauge update with int",
			method:       http.MethodPost,
			path:         "/update/Gauge/Name/42",
			expectedCode: http.StatusOK,
		},
	}

	r := New(&moqMetricRepository{})
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+tt.path, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
