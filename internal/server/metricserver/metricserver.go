package metricserver

import (
	"fmt"
	"net/http"

	"github.com/shurikeagle/metrics-collector/internal/server/metrichandler"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

type metricserver struct {
	server  *http.Server
	handler *metrichandler.Handler
}

func New(ip string, port uint16, storage storage.MetricRepository) *metricserver {
	mServer := &metricserver{
		handler: metrichandler.New(storage),
	}

	mServer.buildHttp(ip, port)

	return mServer
}

func (s *metricserver) Run() error {
	return s.server.ListenAndServe()
}

func (s *metricserver) buildHttp(ip string, port uint16) {
	mux := http.NewServeMux()
	mux.Handle("/update/gauge/", http.HandlerFunc(s.handleUpdateGauge))
	mux.Handle("/update/counter/", http.HandlerFunc(s.handleUpdateCounter))

	addr := fmt.Sprintf("%s:%d", ip, port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(mux.ServeHTTP),
	}
}

func (s *metricserver) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement

	// 0. Decompose path (by split?)
	// 1. Check path (if have existing metric type)
	// 2. Check metric value (if conver to int/float)

	w.WriteHeader(http.StatusOK)
}

func (s *metricserver) handleUpdateGauge(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement

	// 1. Check path (if have existing metric type)
	// 2. Check metric value (if convert to float)

	w.WriteHeader(http.StatusOK)
}
func (s *metricserver) handleUpdateCounter(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement

	// 1. Split path
	// 2. Check metric value (if conver to int)
	// 3.

	w.WriteHeader(http.StatusOK)
}
