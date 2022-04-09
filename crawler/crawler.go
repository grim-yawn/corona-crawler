package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"time"
)

var ErrNoMorePages = errors.New("no more pages")

type Crawler struct {
	*resty.Client
	db *gorm.DB

	// TODO: Find better way to track initial page
	host      string
	tenant    Tenant
	category  Category
	startDate time.Time

	// TODO: Find better way to track last page
	endDate time.Time
}

// NewCrawler which iterates over categoryHistory and stores current page in database
func NewCrawler(db *gorm.DB, host string, tenant Tenant, category Category, startDate, endDate time.Time) *Crawler {
	return &Crawler{Client: resty.New(), db: db, tenant: tenant, category: category, host: host, endDate: endDate, startDate: startDate}
}

func (c *Crawler) NextPage() (*ArticlesPage, error) {
	// Get next page or use default page
	var nextPage PageModel
	result := c.db.Model(&nextPage).Find(&nextPage, "tenant = ? AND category = ?", c.tenant, c.category)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get previous page: %w", result.Error)
	}
	// If no next page for this tenant and category
	if result.RowsAffected == 0 {
		// TODO: Not the best idea but should be good enough for now
		// https://<CMS_API>/<tenant>/categoryHistory/<category>/YYYY/MM/DD
		nextPage.NextPage = fmt.Sprintf("%s/%d/categoryHistory/%d/%s", c.host, c.tenant, c.category, c.startDate.Format("2006/01/02"))
	}

	// TODO: Proper way to parse path params
	// Check if this page is last
	var nextPageDate string
	_, err := fmt.Sscanf(nextPage.NextPage, fmt.Sprintf("%s/%d/categoryHistory/%d/%%s", c.host, c.tenant, c.category), &nextPageDate)
	if err != nil {
		return nil, fmt.Errorf("malformed nextPageUrl %q: %w", nextPage.NextPage, err)
	}
	d, err := time.Parse("2006/01/02", nextPageDate)
	if err != nil {
		return nil, fmt.Errorf("malformed date: %q: %w", nextPageDate, err)
	}

	// TODO: Check if it's actually true
	// If nextPage is before endDate then we already crawled it
	if d.Before(c.endDate) {
		return nil, ErrNoMorePages
	}

	resp, err := c.R().Get(nextPage.NextPage)
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

	// Save next page to db
	err = c.db.Save(&PageModel{
		Tenant:   c.tenant,
		Category: c.category,
		NextPage: *page.Content.NextPage,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save current page to db: %w", err)
	}

	return page, err
}
