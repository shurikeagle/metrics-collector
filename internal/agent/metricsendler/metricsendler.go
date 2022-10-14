package metricsendler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/agent/metric"
	"github.com/shurikeagle/metrics-collector/internal/dto"
)

const sendTimeout = 5 * time.Second
const maxParallelRequests = 10

type sendler struct {
	serverURL      string
	client         *http.Client
	reportInterval time.Duration
}

// New creates new sendler.
// sendler send metrics to configured host with reportInterval
func New(serverAddress string, reportInterval time.Duration) (*sendler, error) {
	if serverAddress == "" {
		return nil, errors.New("address is empty")
	}

	if reportInterval == 0 {
		return nil, errors.New("report can't be 0")
	}

	serverURL := fmt.Sprintf("http://%s", serverAddress)
	if _, err := url.Parse(serverURL); err != nil {
		return nil, err
	}

	return &sendler{
		serverURL:      serverURL,
		reportInterval: reportInterval,
		client:         &http.Client{},
	}, nil
}

func (s *sendler) Run(ctx context.Context, getMetricsFunc func() *metric.Metrics) {
	ticker := time.NewTicker(s.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics := getMetricsFunc()
			s.send(metrics)
		case <-ctx.Done():
			log.Println(ctx.Err(), ", stopping metrics sendler")
			return
		}
	}
}

func (s *sendler) send(metrics *metric.Metrics) {
	sem := make(chan interface{}, maxParallelRequests)

	c := metrics.Counters()
	for m, v := range c {
		metric := dto.Metric{
			ID:    m,
			MType: "Counter",
			Delta: &v,
		}
		go s.makeSendMetricRequest(sem, metric)
	}

	g := metrics.Gauges()
	for m, v := range g {
		metric := dto.Metric{
			ID:    m,
			MType: "Gauge",
			Value: &v,
		}
		go s.makeSendMetricRequest(sem, metric)
	}
}

func (s *sendler) makeSendMetricRequest(sem chan interface{}, metric dto.Metric) {
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), sendTimeout)

	reqBody, err := json.Marshal(metric)
	if err != nil {
		// we send same requests in send func, so err in making request is fatal
		log.Fatal(err)
	}

	mURL := fmt.Sprintf("%s/update", s.serverURL)

	request, err := http.NewRequestWithContext(timeoutCtx, "POST", mURL, bytes.NewBuffer(reqBody))
	if err != nil {
		// we send same requests in send func, so err in making request is fatal
		log.Fatal(err)
	}
	defer cancelFunc()

	request.Header.Add("Content-Type", "application/json")

	sem <- nil
	defer func() { <-sem }()

	response, err := s.client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStr := "response status for" + mURL + "is" + response.Status
		log.Println(errStr)
	}
}
