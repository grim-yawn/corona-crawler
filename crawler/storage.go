package crawler

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Article struct {
	ID        int
	Published time.Time

	// For debug purposes
	Title       string
	TitleHeader string `json:"titleHeader"`
	Lead        string `json:"lead"`
}

func (a Article) toCSV() []string {
	return []string{strconv.Itoa(a.ID), a.Published.String(), a.Title}
}

type ArticleStorage struct {
	f *os.File
	w *csv.Writer
}

func NewArticleStorage(path string) (*ArticleStorage, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	return &ArticleStorage{
		f: f,
		w: csv.NewWriter(f),
	}, nil
}

func (s *ArticleStorage) StoreArticle(a Article) error {
	err := s.w.Write(a.toCSV())
	if err != nil {
		return fmt.Errorf("failed to write to csv file: %w", err)
	}
	return nil
}

func (s *ArticleStorage) Close() error {
	return s.f.Close()
}
