package metricserver

import (
	"context"
	"net/http"

	"github.com/shurikeagle/metrics-collector/internal/server/metrichandler"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

type metricServer struct {
	server *http.Server
}

// New create http metricServer with metric api handling
func New(serverAddress string, storage storage.MetricRepository) *metricServer {
	router := metrichandler.New(storage)
	httpServer := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	return &metricServer{
		server: httpServer,
	}
}

// Run starts metricServer
func (s *metricServer) Run() error {
	return s.server.ListenAndServe()
}

// Run starts metricServer
func (s *metricServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
