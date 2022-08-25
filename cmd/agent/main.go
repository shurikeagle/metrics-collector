package main

import (
	"context"
	"log"
	"os"
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

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	worker.Run(ctx)
}
