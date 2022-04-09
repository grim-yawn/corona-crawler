package crawler

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
	"time"
)

// TODO: Move to env variables (they are not real secrets here, just local dev config)
const cmsHost = "https://feed-prod.unitycms.io"
const postgresDSN = "postgres://corona-user:corona-password@postgres/corona-crawler"

func Run() error {
	db, err := NewDB(postgresDSN)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	//today := time.Now()
	//yearAgo := today.AddDate(-1, 0, 0)

	c := NewCrawler(cmsHost, TenantTwo, CategorySchweiz)
	for range time.Tick(200 * time.Millisecond) {
		page, err := c.NextPage()
		if err != nil {
			return fmt.Errorf("failed to get next page: %w", err)
		}

		models := make([]ArticleModel, 0, len(page.Content.Elements))
		for _, el := range page.Content.Elements {
			if el.IsAboutCovid() {
				models = append(models, ArticleModel{ID: el.ID, Published: el.Content.Published})
			}
		}
		// Skip if page doesn't have any articles about covid
		if len(models) == 0 {
			continue
		}

		// TODO: Should properly retry to ensure that everything saved correctly
		err = db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&models).Error
		if err != nil {
			log.Error().Err(err).Msg("failed to save page")
			continue
		}

		// TODO: Debug
		var count int64
		err = db.Model(&ArticleModel{}).Count(&count).Error
		if err != nil {
			log.Error().Err(err).Msg("failed to count articles")
			continue
		}

		log.Info().Int64("articles", count).Send()
	}

	return nil
}
