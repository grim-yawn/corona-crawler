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
	client *resty.Client
	db     *gorm.DB

	baseURL string

	// TODO: Find better way to track last page
	startDate time.Time
	stopDate  time.Time
}

// TODO: Proper way to parse path params
func getDateFromPageURL(baseURL, pageURL string) (*time.Time, error) {
	var nextPageDateStr string
	_, err := fmt.Sscanf(pageURL, baseURL+"/"+"%s", &nextPageDateStr)
	if err != nil {
		return nil, fmt.Errorf("malformed nextPageUrl %q: %w", pageURL, err)
	}
	nextPageDate, err := time.Parse("2006/01/02", nextPageDateStr)
	if err != nil {
		return nil, fmt.Errorf("malformed date: %q: %w", nextPageDateStr, err)
	}

	return &nextPageDate, nil
}

func getPageURLFromDate(baseURL string, date time.Time) string {
	return fmt.Sprintf("%s/%s", baseURL, date.Format("2006/01/02"))
}

// NewCrawler which iterates over categoryHistory and stores current page in database
func NewCrawler(db *gorm.DB, baseURL string, startDate, endDate time.Time) *Crawler {
	return &Crawler{
		db:      db,
		client:  resty.New(),
		baseURL: baseURL,

		startDate: startDate,
		stopDate:  endDate,
	}
}

func (c *Crawler) NextPage() (*ArticlesPage, error) {
	startPageURL := getPageURLFromDate(c.baseURL, c.startDate)

	// Get old page or set default
	nextPage := &PageModel{}
	result := c.db.Find(&nextPage, "start_page = ?", getPageURLFromDate(c.baseURL, c.startDate))
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get previous page: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		nextPage = &PageModel{
			StartPage: startPageURL,
			NextPage:  startPageURL,
		}
	}

	// Check if we need to stop
	nextPageDate, err := getDateFromPageURL(c.baseURL, nextPage.NextPage)
	if err != nil {
		return nil, err
	}
	if nextPageDate.Before(c.stopDate) {
		return nil, ErrNoMorePages
	}

	// Get next page and move next page in db
	resp, err := c.client.R().Get(nextPage.NextPage)
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
		StartPage: startPageURL,
		NextPage:  *page.Content.NextPage,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save current page to db: %w", err)
	}

	return page, err
}
