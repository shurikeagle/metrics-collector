package metricserver

import (
	"fmt"
	"net/http"

	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

type metricserver struct {
	server  *http.Server
	storage storage.MetricRepository
}

func New(ip string, port uint16, storage storage.MetricRepository) *metricserver {
	mServer := &metricserver{
		storage: storage,
	}

	mServer.buildHttp(ip, port)

	return mServer
}

func (s *metricserver) Run() error {
	return s.server.ListenAndServe()
}

func (s *metricserver) buildHttp(ip string, port uint16) {
	mux := http.NewServeMux()
	mux.Handle("/update", http.HandlerFunc(s.handleUpdate))

	addr := fmt.Sprintf("%s:%d", ip, port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(mux.ServeHTTP),
	}
}

func (s *metricserver) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement in service update metric
	s.storage.AddOrUpdateCounter(metric.Counter{
		Name:  "tst",
		Value: 42,
	})

	w.WriteHeader(http.StatusOK)
}
