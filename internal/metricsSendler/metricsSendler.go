package metricsSendler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/metric"
)

const sendTimeout = 5 * time.Second

type sendler struct {
	serverUrl      string
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

	sUrl := fmt.Sprintf("%s:%d", ip, port)
	if _, err := url.Parse(sUrl); err != nil {
		return nil, err
	}

	return &sendler{
		serverUrl:      sUrl,
		reportInterval: reportInterval,
		client:         &http.Client{},
	}, nil
}

func (s *sendler) Run(ctx context.Context, getMetricsFunc func() metric.Metrics) {
	ticker := time.NewTicker(s.reportInterval)

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

func (s *sendler) send(metrics metric.Metrics) {
	// TODO: Implement max parallel requests logic

	for m, v := range metrics.Counters {
		go s.makeSendMetricRequest("Counter", m, strconv.FormatInt(v, 10))
	}

	for m, v := range metrics.Gauges {
		go s.makeSendMetricRequest("Gauge", m, strconv.FormatFloat(v, 'f', 4, 64))
	}
}

func (s *sendler) makeSendMetricRequest(metricType string, metricName string, value string) {
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), sendTimeout)

	// TODO: Url validation and special symbols handling
	mUrl := fmt.Sprintf("%s/update/%s/%s/%s", s.serverUrl, metricType, metricName, value)

	request, err := http.NewRequestWithContext(timeoutCtx, "POST", mUrl, nil)
	if err != nil {
		// we send same requests in send func, so err in making request is fatal
		log.Fatal(err)
	}
	defer cancelFunc()

	request.Header.Add("Content-Type", "Content-Type: text/plain")

	response, err := s.client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}

	if response.StatusCode != http.StatusOK {
		errStr := "response status for" + mUrl + "is" + response.Status
		log.Println(errStr)
	}
}
