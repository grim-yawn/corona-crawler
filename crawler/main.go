package crawler

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

const cmsHost = "https://feed-prod.unitycms.io"

func Run() error {
	c := NewCrawler(cmsHost, TenantTwo, CategorySchweiz)
	for range time.Tick(100 * time.Millisecond) {
		page, err := c.NextPage()
		if err != nil {
			return fmt.Errorf("failed to get next page: %w", err)
		}

		for _, article := range page.Content.Elements {
			if article.IsAboutCovid() {
				log.Info().Int("article_id", int(article.ID)).Str("title_short", article.Content.Title).Send()
			}
		}
	}

	return nil
}
