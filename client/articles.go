package client

import "time"

type CategoryHistoryResponse struct {
	Content struct {
		NextPage *string `json:"nextpage"`

		Elements []struct {
			ID int `json:"id"`

			Content struct {
				Title       string `json:"title"`
				TitleHeader string `json:"titleHeader"`
				Lead        string `json:"lead"`

				Published time.Time `json:"published"`
			}
		} `json:"elements"`
	}
}
