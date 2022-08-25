package metricsSendler

import (
	"errors"
	"time"
)

type hostAddr struct {
	ip   string
	port uint16
}

type sendler struct {
	host           hostAddr
	reportInterval time.Duration
}

// New creates new sendler.
// sendler send metrics to configured host with reportInterval
func New(ip string, port uint16, reportInterval time.Duration) (*sendler, error) {
	// TODO: add default values to main
	if ip == "" {
		return nil, errors.New("ip is empty")
	}

	if reportInterval == 0 {
		return nil, errors.New("report interval is empty")
	}

	// TODO: Add http client ti sendler

	return &sendler{
		host: hostAddr{
			ip:   ip,
			port: port,
		},
		reportInterval: reportInterval,
	}, nil
}
