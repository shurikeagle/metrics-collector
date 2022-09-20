package metrichandler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/shurikeagle/metrics-collector/internal/server/metric"
)

// POST /update/{metricType}/{metricName}/{metricValue}
func (h *handler) updateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/plain")

	if err := h.updateMetricFromRequest(r); err != nil {
		if err == ErrUnexpectedMetricType {
			w.WriteHeader(http.StatusNotImplemented)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		fmt.Fprintln(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func (h *handler) updateMetricFromRequest(r *http.Request) error {
	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	switch metricType {
	case "gauge":
		if gauge, err := parseGauge(metricName, metricValue); err != nil {
			return err
		} else {
			h.storage.AddOrUpdateGauge(gauge)
		}

	case "counter":
		if counter, err := parseCounter(metricName, metricValue); err != nil {
			return err
		} else {
			h.updateCounter(counter)
		}
	default:
		return ErrUnexpectedMetricType
	}

	return nil
}

func parseGauge(name string, rawValue string) (metric.Gauge, error) {
	if value, err := strconv.ParseFloat(rawValue, 64); err != nil {
		return metric.Gauge{}, err
	} else {
		return metric.Gauge{
			Name:  name,
			Value: value,
		}, nil
	}
}

func parseCounter(name string, rawValue string) (metric.Counter, error) {
	if value, err := strconv.ParseInt(rawValue, 10, 64); err != nil {
		return metric.Counter{}, err
	} else {
		return metric.Counter{
			Name:  name,
			Value: value,
		}, nil
	}
}

// Update updates counter metric in storage
func (h *handler) updateCounter(c metric.Counter) {
	existingCounter, _ := h.storage.GetCounter(c.Name)
	existingCounter.Value += c.Value

	h.storage.AddOrUpdateCounter(existingCounter)
}
