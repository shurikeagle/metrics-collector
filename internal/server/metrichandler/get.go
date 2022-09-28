package metrichandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/shurikeagle/metrics-collector/internal/dto"
)

const metricListTemplate = `<!DOCTYPE html>
<html>
<body>
	<ul> 
		%s
	</ul>
</body>
</html>`

// POST /value
func (h *handler) getValueFromBodyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", JSONcontentType)

	if metric, err := h.getMetricFromBody(r); err != nil {
		errCode := getStatusCodeByError(err)
		w.WriteHeader(errCode)

		errResponse := dto.ErrorResponse{
			Error: err.Error(),
		}
		errResponseBytes, err := json.Marshal(errResponse)
		if err != nil {
			panic(err)
		}

		w.Write(errResponseBytes)
	} else {
		metricBytes, err := json.Marshal(metric)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(metricBytes)
	}
}

// GET /value/{metricType}/{metricName}
func (h *handler) getValueFromPathHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/plain")

	if value, err := h.getMetricValueFromPath(r); err != nil {
		errCode := getStatusCodeByError(err)
		w.WriteHeader(errCode)

		fmt.Fprintln(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, value)
	}
}

// GET /
func (h *handler) getAllHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/html")

	counters, gauges := h.storage.GetAll()

	sb := strings.Builder{}
	for _, c := range counters {
		li := fmt.Sprintf("		<li>%s: %d</li>\n", c.Name, c.Value)
		sb.WriteString(li)
	}
	for _, g := range gauges {
		li := fmt.Sprintf("		<li>%s: %.3f</li>\n", g.Name, g.Value)
		sb.WriteString(li)
	}

	respHTML := fmt.Sprintf(metricListTemplate, sb.String())

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, respHTML)
}

func (h *handler) getMetricFromBody(r *http.Request) (*dto.Metric, error) {

	contentType := r.Header.Get("Content-type")
	if contentType != JSONcontentType {
		return nil, ErrInvalidRequestBody
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	req := dto.GetMetricRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, ErrInvalidRequestBody
	}

	metric := dto.Metric{
		ID:    req.ID,
		MType: req.MType,
	}

	switch strings.ToLower(req.MType) {
	case "gauge":
		if g, ok := h.storage.GetGauge(req.ID); !ok {
			return nil, ErrMetricNotFound
		} else {
			metric.Value = &g.Value
		}
	case "counter":
		if c, ok := h.storage.GetCounter(req.ID); !ok {
			return nil, ErrMetricNotFound
		} else {
			metric.Delta = &c.Value
		}
	default:
		return nil, ErrUnexpectedMetricType
	}

	return &metric, nil
}

func (h *handler) getMetricValueFromPath(r *http.Request) (string, error) {
	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")

	switch metricType {
	case "gauge":
		if g, ok := h.storage.GetGauge(metricName); !ok {
			return "", ErrMetricNotFound
		} else {
			return fmt.Sprintf("%.3f", g.Value), nil
		}
	case "counter":
		if c, ok := h.storage.GetCounter(metricName); !ok {
			return "", ErrMetricNotFound
		} else {
			return fmt.Sprintf("%d", c.Value), nil
		}
	default:
		return "", ErrUnexpectedMetricType
	}
}

func getStatusCodeByError(err error) int {
	switch err {
	case ErrInvalidRequestBody:
		return http.StatusBadRequest
	case ErrUnexpectedMetricType:
		return http.StatusNotImplemented
	case ErrMetricNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
