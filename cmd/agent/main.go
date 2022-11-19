package main

import (
	"context"
	"flag"
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
	ServerAddress  string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

var cfg *appConfig = &appConfig{}

func init() {
	flag.StringVar(&cfg.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "Agent poller's poll interval")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "Agent report interval to server")
}

func main() {
	log.Println("poll agent start")

	buildAppConfig()

	rPoller := runtimepoller.Poller{}
	worker, err := pollworker.New(&rPoller, cfg.PollInterval)
	if err != nil {
		log.Fatal(err)
	}

	mSedler, err := metricsendler.New(cfg.ServerAddress, cfg.ReportInterval)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		if err := worker.Run(ctx); err != nil {
			log.Println(err)
		}
	}()
	go mSedler.Run(ctx, worker.Stats)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("agent stopped")
}

func buildAppConfig() {
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
