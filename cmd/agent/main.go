package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/agent/metricsendler"
	"github.com/shurikeagle/metrics-collector/internal/agent/pollworker"
	"github.com/shurikeagle/metrics-collector/internal/agent/runtimepoller"
)

const (
	serverIP   = "http://127.0.0.1"
	serverPort = 8080
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	log.Println("poll agent start")

	rPoller := runtimepoller.Poller{}
	worker, err := pollworker.New(&rPoller, pollInterval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	mSedler, err := metricsendler.New(serverIP, serverPort, reportInterval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go worker.Run(ctx)
	go mSedler.Run(ctx, worker.Stats)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("agent stopped")
}
