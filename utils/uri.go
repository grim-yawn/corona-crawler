package utils

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

func GetCategoryHistoryURLFromDate(baseURL string, category int, date time.Time) string {
	return fmt.Sprintf("%s/categoryHistory/%d/%s", baseURL, category, date.Format("2006/01/02"))
}

func GetDateFromCategoryHistoryURL(uri string, category int) (*time.Time, error) {
	// Format <HOST>/<tenant>/categoryHistory/<category>/<date>
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid url %q: %v", uri, err)
	}

	// TODO: Use something more reliable
	parts := strings.SplitN(u.Path, "/", 7)
	if len(parts) != 7 {
		return nil, fmt.Errorf("malformed path %q: %w", u.Path, err)
	}
	dateStr := strings.Join(parts[4:], "/")
	d, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date part of url %q: %w", dateStr, err)
	}

	return &d, nil
}
