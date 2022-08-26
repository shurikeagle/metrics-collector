package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/pollWorker"
	"github.com/shurikeagle/metrics-collector/internal/runtimePoller"
)

const (
	serverIp   = "http://localhost"
	serverPort = "8080"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	rPoller := runtimePoller.Poller{}
	worker, err := pollWorker.New(2*time.Second, &rPoller)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go worker.Run(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-t.C:
			// TODO: Debug, to remove
			curStats := worker.Stats()
			log.Println("============================")
			for k, v := range curStats.Gauges {
				log.Println(k, ":", v)
			}
			for k, v := range curStats.Counters {
				log.Println(k, ":", v)
			}
			log.Println("============================")
		case <-quit:
			log.Println("agent stopped")
			os.Exit(0)
		}
	}

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// <-quit
	// log.Println("agent stopped")
}
