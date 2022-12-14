package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shurikeagle/metrics-collector/internal/server/metricserver"
	"github.com/shurikeagle/metrics-collector/internal/server/storage/inmemory"
)

const (
	serverIP   = "127.0.0.1"
	serverPort = 8080
)

func main() {
	log.Println("starting metric server")
	storage := inmemory.New()
	mServer := metricserver.New(
		serverIP,
		serverPort,
		storage,
	)

	go func() {
		log.Fatal(mServer.Run())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("metric server stopped")
}
