package latest

import (
	"corona-crawler/articles"
	"corona-crawler/client"
)

type Crawler struct {
	client   *client.Client
	category int

	matcher articles.Matcher
}

// NewCrawler iterates over latest articles in "category" nonstop
func NewCrawler(articlesClient *client.Client, category int, matcher articles.Matcher) *Crawler {
	return &Crawler{
		client:   articlesClient,
		category: category,
		matcher:  matcher,
	}
}

// NextPage returns slice of articles.Article
// Result can be an empty slice and error=nil, if all pages were filtered
func (c *Crawler) NextPage() ([]articles.Article, error) {
	category, err := c.client.GetCategory(c.category)
	if err != nil {
		return nil, err
	}

	// Filter articles
	filteredArticles := make([]articles.Article, 0, len(category.Content.Elements.Articles))
	for _, article := range category.Content.Elements.Articles {
		if !c.matcher.MatchArticle(article) {
			continue
		}
		filteredArticles = append(filteredArticles, article)
	}

	return filteredArticles, nil
}
