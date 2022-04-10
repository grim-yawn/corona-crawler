package crawler

import (
	"errors"
	"fmt"
	"time"

	"corona-crawler/articles"
	"corona-crawler/client"
	"corona-crawler/utils"
)

var ErrNoMorePages = errors.New("no more pages")

type Crawler struct {
	client *client.Client

	stateStorage *StateStorage

	// TODO: Find better way to track last page
	category  int
	startDate time.Time
	stopDate  time.Time

	matcher articles.Matcher
}

// NewCrawler which iterates over categoryHistory and stores current page in database
func NewCrawler(
	stateStorage *StateStorage,
	articlesClient *client.Client,
	category int, startDate time.Time, endDate time.Time,
	matcher articles.Matcher,
) *Crawler {
	return &Crawler{
		stateStorage: stateStorage,
		client:       articlesClient,

		category: category,

		startDate: startDate,
		stopDate:  endDate,

		matcher: matcher,
	}
}

// NextPage returns slice of ArticleModels or ErrNoMorePages
// Result can be an empty slice and error=nil, if all pages were filtered
func (c *Crawler) NextPage() ([]articles.Article, error) {
	startPageURL := utils.GetCategoryHistoryURLFromDate(c.client.BaseURL, c.category, c.startDate)

	// Get old page or set default
	state, err := c.stateStorage.GetStateOrDefault(startPageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from db: %w", err)
	}

	// Check if we need to stop
	nextPageDate, err := utils.GetDateFromCategoryHistoryURL(state.NextPage, c.category)
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
	filteredArticles := make([]articles.Article, 0, len(history.Content.Elements))
	for _, article := range history.Content.Elements {
		if !c.matcher.MatchArticle(article) {
			continue
		}
		filteredArticles = append(filteredArticles, article)
	}

	// Save next page to db
	// TODO: Potential data loss bc page set before data saved to database
	state.NextPage = *history.Content.NextPage
	err = c.stateStorage.SaveState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to save current page to db: %w", err)
	}

	return filteredArticles, nil
}
