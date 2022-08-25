package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shurikeagle/metrics-collector/internal/collectWorker"
	"github.com/shurikeagle/metrics-collector/internal/runtimeCollector"
)

func main() {
	rCollector := runtimeCollector.Collector{}
	worker, err := collectWorker.New(2*time.Second, &rCollector)
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
