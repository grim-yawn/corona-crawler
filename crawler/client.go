package crawler

import (
	"fmt"
	"time"
)

// Category specifies which category we need to crawl
type Category int

const (
	Schweiz Category = 6
)

// Tenant part of request
// TODO: Need to understand what does it mean
// Using just tenant=2 like in browser request
// Testing different values, available range seems to be tenant=[1...15]
// Content a bit different and language also alternates between DE and FR
type Tenant int

const (
	TenantTwo = 2
)

type CategoryHistoryRequest struct {
	Host string

	Tenant   Tenant
	Category Category
}

// String returns string url for this category and tenant
// TODO: Sprintf not the best solution but should be enough, bc this will be used only in one request
func (r CategoryHistoryRequest) String() string {
	return fmt.Sprintf("%s/%d/categoryHistory/%d", r.Host, r.Tenant, r.Category)
}

// CategoryHistoryResponse is collection of all required fields from full contentHistory response
// TODO: Lots of nested structure definitions
type CategoryHistoryResponse struct {
	// Current server date
	Date time.Time `json:"date"`

	Content struct {
		NextPage string `json:"nextpage"`
		Elements []struct {
			ID      int `json:"id"`
			Content struct {
				Title       string    `json:"title"`
				TitleHeader string    `json:"titleHeader"`
				Lead        string    `json:"lead"`
				Published   time.Time `json:"published"`

				// Just in case we need to track updates aswell
				Updated time.Time `json:"updated"`
			} `json:"content"`
		} `json:"elements"`
	} `json:"content"`
}
