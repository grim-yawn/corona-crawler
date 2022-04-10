package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"corona-crawler/articles"
)

type CategoryHistoryResponse struct {
	Content struct {
		NextPage *string `json:"nextpage"`

		Elements []articles.Article `json:"elements"`
	} `json:"content"`
}

type CategoryResponse struct {
	Content struct {
		Elements CategoryResponseElements `json:"elements"`
	} `json:"content"`
}

type CategoryResponseElements struct {
	CategoryHistory string
	Articles        []articles.Article
}

func (e *CategoryResponseElements) UnmarshalJSON(data []byte) error {
	rawJSON := make([]map[string]interface{}, 0)
	err := json.Unmarshal(data, &rawJSON)
	if err != nil {
		return err
	}

	if len(rawJSON) == 0 {
		return errors.New("can't be empty")
	}

	// Category history
	history, exists := rawJSON[0]["categoryHistory"]
	if !exists {
		return errors.New("first element must be history")
	}
	historyStr, ok := history.(string)
	if !ok {
		return fmt.Errorf("history must be string, got %T instead", history)
	}
	e.CategoryHistory = historyStr

	// Remaining elements
	// TODO: Should use smth better, don't want to parse maps here
	articlesData, err := json.Marshal(rawJSON[1:])
	if err != nil {
		return err
	}
	e.Articles = make([]articles.Article, len(rawJSON)-1)
	err = json.Unmarshal(articlesData, &e.Articles)
	if err != nil {
		return err
	}

	return nil
}
