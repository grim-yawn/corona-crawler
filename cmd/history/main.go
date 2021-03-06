package main

import (
	"time"

	"github.com/rs/zerolog/log"

	"corona-crawler/articles"
	"corona-crawler/client"
	"corona-crawler/db"
	"corona-crawler/history"
)

// Category specifies which category we need to crawl
type Category = int

const (
	CategorySchweiz Category = 6
)

// Tenant dummy type for document
// Accepts [1..15] return different articles, language alternate between DE and FR
type Tenant = int

// TODO: Move to env variables (they are not real secrets here, just local dev config)
// TODO: https://<HOST>/<tenant>/categoryHistory/<category>, check Tenant and Category
const cmsBaseURL = "https://feed-prod.unitycms.io/2"
const postgresDSN = "postgres://corona-user:corona-password@postgres/corona-crawler"

func main() {
	log.Info().Msg("crawler started")

	gormDB, err := db.New(postgresDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database")
	}

	// TODO: Move to config
	// TODO: Use it for parallel processing: split this date range into N parts and run them in parallel
	// Should work as intended bc current state saved by StartDate key
	today := time.Now()
	yearAgo := today.AddDate(-1, 0, 0)

	articleClient := client.New(cmsBaseURL)
	stateStorage := crawler.NewStateStorage(gormDB)
	articlesStorage := articles.NewStorage(gormDB)

	// TODO: Proper config
	c := crawler.NewCrawler(stateStorage, articleClient, CategorySchweiz, today, yearAgo, articles.ArticleAboutCovid)
	for range time.Tick(200 * time.Millisecond) {
		page, err := c.NextPage()
		if err == crawler.ErrNoMorePages {
			log.Info().Msg("successfully parsed all articles")
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to get next page")
			break
		}
		// Skip if page doesn't have any articles about covid
		if len(page) == 0 {
			continue
		}

		// TODO: Should properly retry to ensure that everything saved correctly
		err = articlesStorage.BatchSaveArticles(articles.ToModels(page))
		if err != nil {
			log.Error().Err(err).Msg("failed to save articles")
			break
		}
	}
}
