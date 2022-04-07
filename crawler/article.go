package crawler

import (
	"strings"
	"time"
)

// Category specifies which category we need to crawl
type Category int

const (
	CategorySchweiz Category = 6
)

// Tenant part of request
// TODO: Need to understand what does it mean
// Using just tenant=2 like in the request from browser
// Testing different values, available range seems to be tenant=[1...15]
// Content a bit different and language also alternates between DE and FR
type Tenant int

const (
	TenantTwo = 2
)

type ArticleID int

type Article struct {
	ID      ArticleID `json:"id"`
	Content struct {
		Title       string `json:"title"`
		TitleHeader string `json:"titleHeader"`
		Lead        string `json:"lead"`

		Published time.Time `json:"published"`
	}
}

// IsAboutCovid checks if this article is about covid-19
func (a Article) IsAboutCovid() bool {
	// TODO: Not the smartest way to compare strings but don't want to use regexp here
	// TODO: Proper case insensitive match with regexp?
	for _, sub := range []string{"Corona", "Covid-19"} {
		if strings.Contains(a.Content.Title, sub) {
			return true
		}
		if strings.Contains(a.Content.TitleHeader, sub) {
			return true
		}
		if strings.Contains(a.Content.Lead, sub) {
			return true
		}
	}

	return false
}

type ArticlesPage struct {
	Content struct {
		Elements []Article `json:"elements"`
		NextPage *string   `json:"nextpage"`
	} `json:"content"`
}
