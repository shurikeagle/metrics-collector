package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/shurikeagle/metrics-collector/internal/server/metricserver"
	"github.com/shurikeagle/metrics-collector/internal/server/storage/inmemory"
)

type appConfig struct {
	ServerAddress           string        `env:"ADDRESS"`
	RepArchiveStoreInterval time.Duration `env:"STORE_INTERVAL"`
	RepArchiveStoreFile     string        `env:"STORE_FILE"`
	RestoreRepOnStart       bool          `env:"RESTORE"`
}

var cfg *appConfig = &appConfig{}

func init() {
	flag.StringVar(&cfg.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.BoolVar(&cfg.RestoreRepOnStart, "r", true, "Restore repository archive on start")
	flag.DurationVar(&cfg.RepArchiveStoreInterval, "i", 300*time.Second, "Repository archive store interval")
	flag.StringVar(&cfg.RepArchiveStoreFile, "f", "/tmp/devops-metrics-db.json", "Repository store file full name")
}

func main() {
	log.Println("starting metric server")

	buildAppConfig()

	inmemArchiveSettings := inmemory.InmemArchiveSettings{
		StoreInterval:   cfg.RepArchiveStoreInterval,
		FileName:        cfg.RepArchiveStoreFile,
		RestoreOnCreate: cfg.RestoreRepOnStart,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	storage := inmemory.New(inmemArchiveSettings, ctx)

	mServer := metricserver.New(cfg.ServerAddress, storage)

	go func() {
		log.Fatal(mServer.Run())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit

	if err := storage.ArchiveAll(); err != nil {
		log.Printf("couldn't archive metrics on stop: %s", err.Error())
	} else {
		log.Printf("metrics were archived")
	}
	log.Println("metric server stopped")
}

func buildAppConfig() {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("before parse: %v\n", cfg)

	flag.Parse()

	fmt.Printf("after parse: %v\n", cfg)
}
