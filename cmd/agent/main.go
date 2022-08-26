package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	metricsSendler "github.com/shurikeagle/metrics-collector/internal/metricsSendler"
	"github.com/shurikeagle/metrics-collector/internal/pollWorker"
	"github.com/shurikeagle/metrics-collector/internal/runtimePoller"
)

const (
	serverIp   = "http://127.0.0.1"
	serverPort = 8080
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	rPoller := runtimePoller.Poller{}
	worker, err := pollWorker.New(&rPoller, pollInterval)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	mSedler, err := metricsSendler.New(serverIp, serverPort, reportInterval)
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
