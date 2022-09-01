package metricserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

type moqMetricRepository struct{}

var _ storage.MetricRepository = (*moqMetricRepository)(nil)

func (r *moqMetricRepository) GetCounter(name string) (c metric.Counter, ok bool) {
	return metric.Counter{}, true
}
func (r *moqMetricRepository) AddOrUpdateGauge(g metric.Gauge)     {}
func (r *moqMetricRepository) AddOrUpdateCounter(c metric.Counter) {}

func Test_metricserver_handleUpdate(t *testing.T) {

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
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "update metric without value",
			method:       http.MethodPost,
			path:         "/update/Gauge/Name/",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "update unexpected metric",
			method:       http.MethodPost,
			path:         "/update/Unexpected/Name/42",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid hhtp method",
			method:       http.MethodGet,
			path:         "/update/Counter/Name/42",
			expectedCode: http.StatusNotFound,
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

	server := New("127.0.0.1", 8080, &moqMetricRepository{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			request.Header.Add("Content-type", "text/plain")

			w := httptest.NewRecorder()
			h := http.HandlerFunc(server.handleUpdate)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.expectedCode, res.StatusCode)
		})
	}
}
