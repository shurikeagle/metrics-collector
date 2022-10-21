package main

import (
	"context"
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
	ServerAddress           string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	RepArchiveStoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	RepArchiveStoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	RestoreRepOnStart       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	log.Println("starting metric server")

	appConfig := buildAppConfig()

	inmemArchiveSettings := inmemory.InmemArchiveSettings{
		StoreInterval:   appConfig.RepArchiveStoreInterval,
		FileName:        appConfig.RepArchiveStoreFile,
		RestoreOnCreate: appConfig.RestoreRepOnStart,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	storage := inmemory.New(inmemArchiveSettings, ctx)

	mServer := metricserver.New(appConfig.ServerAddress, storage)

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

func buildAppConfig() appConfig {
	cfg := appConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
