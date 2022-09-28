package metrichandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/shurikeagle/metrics-collector/internal/dto"
	"github.com/shurikeagle/metrics-collector/internal/server/metric"
)

var ErrEmptyValueForGauge = errors.New("field 'value' cannot be empty for gauge metric")
var ErrEmptyValueForCounter = errors.New("field 'delta' cannot be empty for counter metric")

// TODO: tests for this method
// POST /update
func (h *handler) updateHandlerFromBody(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	if contentType != JSONcontentType {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, JSONcontentTypeExpected)
	} else {
		h.updateMetric(w, r, h.updateMetricFromBody)
	}
}

// POST /update/{metricType}/{metricName}/{metricValue}
func (h *handler) updateHandlerFromPath(w http.ResponseWriter, r *http.Request) {
	h.updateMetric(w, r, h.updateMetricFromPath)
}

func (h *handler) updateMetric(w http.ResponseWriter, r *http.Request, updateFunc func(r *http.Request) error) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/plain")

	if err := updateFunc(r); err != nil {
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

func (h *handler) updateMetricFromBody(r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	updateRequest := dto.Metric{}
	if err = json.Unmarshal(body, &updateRequest); err != nil {
		return err
	}

	switch strings.ToLower(updateRequest.MType) {
	case "gauge":
		if gauge, err := toGauge(updateRequest); err != nil {
			return err
		} else {
			h.storage.AddOrUpdateGauge(*gauge)
		}

	case "counter":
		if counter, err := toCounter(updateRequest); err != nil {
			return err
		} else {
			h.updateCounter(*counter)
		}

	default:
		return ErrUnexpectedMetricType
	}

	return nil
}

func toGauge(m dto.Metric) (*metric.Gauge, error) {
	if m.Value == nil {
		return nil, ErrEmptyValueForGauge
	}

	return &metric.Gauge{
		Name:  m.ID,
		Value: *m.Value,
	}, nil
}

func toCounter(m dto.Metric) (*metric.Counter, error) {
	if m.Delta == nil {
		return nil, ErrEmptyValueForGauge
	}

	return &metric.Counter{
		Name:  m.ID,
		Value: *m.Delta,
	}, nil
}

func (h *handler) updateMetricFromPath(r *http.Request) error {
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
