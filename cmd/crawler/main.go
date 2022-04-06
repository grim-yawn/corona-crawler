package main

import (
	"corona-crawler/crawler"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
)

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter())
	logger.Info().Msg("Crawler started")

	go crawler.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done

	logger.Info().Msg("shutting down")
}
