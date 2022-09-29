package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/shurikeagle/metrics-collector/internal/agent/metricsendler"
	"github.com/shurikeagle/metrics-collector/internal/agent/pollworker"
	"github.com/shurikeagle/metrics-collector/internal/agent/runtimepoller"
)

type appConfig struct {
	ServerAddress     string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	PollIntervalSec   int    `env:"POLL_INTERVAL" envDefault:"2"`
	ReportIntervalSec int    `env:"REPORT_INTERVAL" envDefault:"10"`
}

func main() {
	log.Println("poll agent start")

	appConfig := buildAppConfig()

	rPoller := runtimepoller.Poller{}
	pollInterval := time.Duration(appConfig.PollIntervalSec) * time.Second
	worker, err := pollworker.New(&rPoller, pollInterval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	reportInterval := time.Duration(appConfig.ReportIntervalSec) * time.Second
	mSedler, err := metricsendler.New(appConfig.ServerAddress, reportInterval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		log.Println(worker.Run(ctx))
	}()
	go mSedler.Run(ctx, worker.Stats)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("agent stopped")
}

func buildAppConfig() appConfig {
	cfg := appConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
