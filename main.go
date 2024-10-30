package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

func main() {
	setupLogging()

	ticker := time.NewTicker(5 * time.Second)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go monitorTemperature(ticker, ctx.Done())

	<-ctx.Done()
	log.Infof("received sigterm or sigint, shutting down")
}

func setupLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
}
