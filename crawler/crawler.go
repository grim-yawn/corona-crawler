package crawler

import (
	"corona-crawler/client"
	"corona-crawler/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Article struct {
	ID int

	Title     string
	TitleHead string
	Lead      string

	Published time.Time
}

type ArticleFilterer interface {
	FilterArticle(article Article) bool
}

type ArticleFilterFunc func(article Article) bool

func (f ArticleFilterFunc) FilterArticle(article Article) bool {
	return f(article)
}

func ArticleAboutCovid(article Article) bool {
	// TODO: Not the smartest way to compare strings but don't want to use regexp here
	// TODO: Proper case insensitive match with regexp?
	for _, sub := range []string{"Corona", "Covid-19"} {
		if strings.Contains(article.Title, sub) {
			return true
		}
		if strings.Contains(article.TitleHead, sub) {
			return true
		}
		if strings.Contains(article.Lead, sub) {
			return true
		}
	}

	return false
}

var ErrNoMorePages = errors.New("no more pages")

type Crawler struct {
	client *client.Client
	db     *gorm.DB

	category int

	// TODO: Find better way to track last page
	startDate time.Time
	stopDate  time.Time

	//
	filter ArticleFilterer
}

// NewCrawler which iterates over categoryHistory and stores current page in database
func NewCrawler(db *gorm.DB, c *client.Client, category int, startDate, endDate time.Time, filter ArticleFilterer) *Crawler {
	return &Crawler{
		db:     db,
		client: c,

		category: category,

		startDate: startDate,
		stopDate:  endDate,

		filter: filter,
	}
}

// NextPage returns slice of ArticleModels or ErrNoMorePages
// Result can be an empty slice and error=nil, if all pages were filtered
func (c *Crawler) NextPage() ([]ArticleModel, error) {
	startPageURL := utils.GetCategoryHistoryURLFromDate(c.client.BaseURL, c.category, c.startDate)

	// Get old page or set default
	nextPage := &PageModel{}
	result := c.db.Find(&nextPage, "start_page = ?", startPageURL)
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
	nextPageDate, err := utils.GetDateFromCategoryHistoryURL(nextPage.NextPage, c.category)
	if err != nil {
		return nil, err
	}
	if nextPageDate.Before(c.stopDate) {
		return nil, ErrNoMorePages
	}

	history, err := c.client.GetCategoryHistory(c.category, *nextPageDate)
	if err != nil {
		return nil, err
	}
	if history.Content.NextPage == nil {
		return nil, ErrNoMorePages
	}

	// Filter articles
	articles := make([]ArticleModel, 0, len(history.Content.Elements))
	for _, el := range history.Content.Elements {
		article := Article{
			ID:        el.ID,
			Title:     el.Content.Title,
			TitleHead: el.Content.TitleHeader,
			Lead:      el.Content.Lead,
			Published: el.Content.Published,
		}

		if !c.filter.FilterArticle(article) {
			continue
		}

		articles = append(articles, ArticleModel{ID: article.ID, Published: el.Content.Published})
	}

	// Save next page to db
	err = c.db.Save(&PageModel{
		StartPage: startPageURL,
		NextPage:  *history.Content.NextPage,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save current page to db: %w", err)
	}

	return articles, nil
}
