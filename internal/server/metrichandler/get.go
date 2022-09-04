package metrichandler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

const metricListTemplate = `<!DOCTYPE html>
<html>
<body>
	<ul> 
		%s
	</ul>
</body>
</html>`

// GET /value/{metricType}/{metricName}
func (h *handler) getValueHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "text/plain")

		if value, err := h.getMetricValue(r); err != nil {
			switch err {
			case ErrUnexpectedMetricType:
				w.WriteHeader(http.StatusNotImplemented)
			case ErrMetricNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			fmt.Fprintln(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, value)
		}
	}
}

// GET /
func (h *handler) getAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "text/html")

		counters, gauges := h.storage.GetAll()

		sb := strings.Builder{}
		for _, c := range counters {
			li := fmt.Sprintf("		<li>%s: %d</li>\n", c.Name, c.Value)
			sb.WriteString(li)
		}
		for _, g := range gauges {
			li := fmt.Sprintf("		<li>%s: %f</li>\n", g.Name, g.Value)
			sb.WriteString(li)
		}

		respHtml := fmt.Sprintf(metricListTemplate, sb.String())

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, respHtml)
	}
}

func (h *handler) getMetricValue(r *http.Request) (string, error) {
	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")

	switch metricType {
	case "gauge":
		if g, ok := h.storage.GetGauge(metricName); !ok {
			return "", ErrMetricNotFound
		} else {
			return fmt.Sprintf("%f", g.Value), nil
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
