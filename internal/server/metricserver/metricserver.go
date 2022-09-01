package metricserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/shurikeagle/metrics-collector/internal/server/metric"
	"github.com/shurikeagle/metrics-collector/internal/server/metrichandler"
	"github.com/shurikeagle/metrics-collector/internal/server/storage"
)

var ErrInvalidPathFormat = errors.New("relative path should be `udpate/metricType/metricName/metricValue`")
var ErrUnexpectedMetricType = errors.New("unexpected metric type")

type decomposedUpdatePath struct {
	metricType  string
	metricName  string
	metricValue string
}

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
	mux.Handle("/update/", http.HandlerFunc(s.handleUpdate))

	addr := fmt.Sprintf("%s:%d", ip, port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(mux.ServeHTTP),
	}
}

func (s *metricserver) handleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "invalid http method")

		return
	}

	contentType := r.Header.Get("Content-type")
	if contentType != "text/plain" && contentType != "text/plain; charset=utf-8" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid content type", contentType)

		return
	}

	decomposedPath, err := decomposeUpdateMetricPath(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)

		return
	}

	if err := s.updateMetricByDecomposedPath(decomposedPath); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func decomposeUpdateMetricPath(path string) (decomposedUpdatePath, error) {
	path = strings.Trim(path, "/")
	splited := strings.Split(path, "/")
	if len(splited) != 4 {
		return decomposedUpdatePath{}, ErrInvalidPathFormat
	}

	return decomposedUpdatePath{
		metricType:  splited[1],
		metricName:  splited[2],
		metricValue: splited[3],
	}, nil
}

func (s *metricserver) updateMetricByDecomposedPath(decomposedPath decomposedUpdatePath) error {
	metricType := strings.ToLower(decomposedPath.metricType)
	switch metricType {
	case "gauge":
		if gauge, err := getGauge(decomposedPath.metricName, decomposedPath.metricValue); err != nil {
			return err
		} else {
			s.handler.UpdateGauge(gauge)
		}

	case "counter":
		if counter, err := getCounter(decomposedPath.metricName, decomposedPath.metricValue); err != nil {
			return err
		} else {
			s.handler.UpdateCounter(counter)
		}
	default:
		return ErrUnexpectedMetricType
	}

	return nil
}

func getGauge(name string, rawValue string) (metric.Gauge, error) {
	if value, err := strconv.ParseFloat(rawValue, 64); err != nil {
		return metric.Gauge{}, err
	} else {
		return metric.Gauge{
			Name:  name,
			Value: value,
		}, nil
	}
}

func getCounter(name string, rawValue string) (metric.Counter, error) {
	if value, err := strconv.ParseInt(rawValue, 10, 64); err != nil {
		return metric.Counter{}, err
	} else {
		return metric.Counter{
			Name:  name,
			Value: value,
		}, nil
	}
}
