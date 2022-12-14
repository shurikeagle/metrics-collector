package metricsendler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/agent/metric"
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
func New(ip string, port uint16, reportInterval time.Duration) (*sendler, error) {
	if ip == "" {
		return nil, errors.New("ip is empty")
	}

	if reportInterval == 0 {
		return nil, errors.New("report can't be 0")
	}

	sURL := fmt.Sprintf("%s:%d", ip, port)
	if _, err := url.Parse(sURL); err != nil {
		return nil, err
	}

	return &sendler{
		serverURL:      sURL,
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
		go s.makeSendMetricRequest(sem, "Counter", m, strconv.FormatInt(v, 10))
	}

	g := metrics.Gauges()
	for m, v := range g {
		go s.makeSendMetricRequest(sem, "Gauge", m, strconv.FormatFloat(v, 'f', 4, 64))
	}
}

func (s *sendler) makeSendMetricRequest(sem chan interface{}, metricType string, metricName string, value string) {
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), sendTimeout)

	// TODO: Check if need url validation and special symbols handling
	mURL := fmt.Sprintf("%s/update/%s/%s/%s", s.serverURL, metricType, metricName, value)

	request, err := http.NewRequestWithContext(timeoutCtx, "POST", mURL, nil)
	if err != nil {
		// we send same requests in send func, so err in making request is fatal
		log.Fatal(err)
	}
	defer cancelFunc()

	request.Header.Add("Content-Type", "text/plain")

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
