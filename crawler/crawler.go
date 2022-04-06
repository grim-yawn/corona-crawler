package crawler

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"time"
)

func Run() {
	storage, err := NewArticleStorage("./output.csv")
	if err != nil {
		log.Error().Err(err).Msg("failed to create articles storage")
		return
	}
	defer func() {
		err := storage.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close storage")
		}
	}()

	client := resty.New()

	nextItemURL := CategoryHistoryRequest{
		Host:     "https://feed-prod.unitycms.io",
		Tenant:   TenantTwo,
		Category: Schweiz,
	}.String()

	for range time.Tick(200 * time.Millisecond) {
		resp, err := client.R().Get(nextItemURL)
		if err != nil {
			log.Error().Err(err).Msg("failed to get data from API")
			return
		}
		if !resp.IsSuccess() {
			log.Error().Int("status_code", resp.StatusCode()).Msg("non success status from api")
			return
		}

		r := &CategoryHistoryResponse{}
		err = json.Unmarshal(resp.Body(), r)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse response")
			return
		}

		for _, el := range r.Content.Elements {
			err = storage.StoreArticle(Article{
				ID:        el.ID,
				Published: el.Content.Published,
				Title:     el.Content.Title,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to save item to storage")
				return
			}
		}
	}
}
