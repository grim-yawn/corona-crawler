package client

import (
	"corona-crawler/utils"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Client struct {
	*resty.Client
}

func New(baseURL string) *Client {
	return &Client{
		Client: resty.New().SetBaseURL(baseURL),
	}
}

// GetCategoryHistory returns archive articles published before endTime
// Article page contains articles for 2-3 full days
func (c *Client) GetCategoryHistory(category int, endDate time.Time) (*CategoryHistoryResponse, error) {
	resp, err := c.R().Get(utils.GetCategoryHistoryURLFromDate(c.BaseURL, category, endDate))
	if err != nil {
		return nil, fmt.Errorf("failed to get data from API: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("bad status from API (%d, %s)", resp.StatusCode(), resp.Status())
	}

	history := &CategoryHistoryResponse{}
	err = json.Unmarshal(resp.Body(), history)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	return history, nil
}

// GetCategory return latest articles in category
func (c *Client) GetCategory(category int) error {
	return nil
}
