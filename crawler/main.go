package crawler

import (
	"corona-crawler/client"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
	"time"
)

// TODO: Move to env variables (they are not real secrets here, just local dev config)
// TODO: https://<HOST>/<tenant>/categoryHistory/<category>, check Tenant and Category
const cmsBaseURL = "https://feed-prod.unitycms.io/2"
const postgresDSN = "postgres://corona-user:corona-password@postgres/corona-crawler"

func Run() error {
	db, err := NewDB(postgresDSN)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	today := time.Now()
	yearAgo := today.AddDate(-1, 0, 0)

	articleClient := client.New(cmsBaseURL)

	c := NewCrawler(db, articleClient, CategorySchweiz, today, yearAgo, ArticleFilterFunc(ArticleAboutCovid))
	for range time.Tick(200 * time.Millisecond) {
		page, err := c.NextPage()
		if err == ErrNoMorePages {
			log.Info().Msg("successfully parse all articles")
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to get next page: %w", err)
		}

		// Skip if page doesn't have any articles about covid
		if len(page) == 0 {
			continue
		}

		// TODO: Should properly retry to ensure that everything saved correctly
		err = db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&page).Error
		if err != nil {
			log.Error().Err(err).Msg("failed to save page")
			continue
		}
	}

	return nil
}
