package main

import (
	"corona-crawler/crawler"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("crawler started")
	err := crawler.Run()
	if err != nil {
		log.Error().Err(err).Msg("crawler stopped after error")
	}
}
