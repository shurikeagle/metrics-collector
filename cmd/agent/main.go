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
	ServerAddress  string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
}

func main() {
	log.Println("poll agent start")

	appConfig := buildAppConfig()

	rPoller := runtimepoller.Poller{}
	worker, err := pollworker.New(&rPoller, appConfig.PollInterval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	mSedler, err := metricsendler.New(appConfig.ServerAddress, appConfig.ReportInterval)
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
