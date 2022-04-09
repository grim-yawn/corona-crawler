package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

var ErrNoMorePages = errors.New("no more pages")

type Crawler struct {
	*resty.Client

	nextPage string
}

// NewCrawler
// TODO: Not the best idea but should be good enough for now
// https://<CMS_API>/<tenant>/categoryHistory/<category>/[YYYY/MM/DD]
func NewCrawler(host string, tenant Tenant, category Category) *Crawler {
	c := &Crawler{Client: resty.New()}

	// First page
	c.nextPage = fmt.Sprintf("%s/%d/categoryHistory/%d", host, tenant, category)

	return c
}

func (c *Crawler) NextPage() (*ArticlesPage, error) {
	resp, err := c.R().Get(c.nextPage)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from API: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("bad status from API (%d, %s)", resp.StatusCode(), resp.Status())
	}

	page := &ArticlesPage{}
	err = json.Unmarshal(resp.Body(), page)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if page.Content.NextPage == nil {
		return nil, ErrNoMorePages
	}

	c.nextPage = *page.Content.NextPage

	return page, err
}
