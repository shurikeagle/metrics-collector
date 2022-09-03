package metricserver

import (
	"fmt"
	"net/http"

	"github.com/shurikeagle/metrics-collector/internal/server/metrichandler"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

type metricServer struct {
	server *http.Server
}

// New create http metricServer with metric api handling
func New(ip string, port uint16, storage storage.MetricRepository) *metricServer {
	addr := fmt.Sprintf("%s:%d", ip, port)
	router := metrichandler.New(storage)
	httpServer := &http.Server{
		Addr:    addr,
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
