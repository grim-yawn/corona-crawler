package main

import (
	"corona-crawler/server"
	"github.com/rs/zerolog/log"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatal().Err(err)
	}
}
