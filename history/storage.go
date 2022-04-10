package crawler

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// State tracks next page for crawling to avoid redundant crawling
type State struct {
	StartPage string `gorm:"primaryKey"`
	NextPage  string
}

func (p State) TableName() string {
	return "crawler_state"
}

type StateStorage struct {
	db *gorm.DB
}

func NewStateStorage(db *gorm.DB) *StateStorage {
	return &StateStorage{db: db}
}

func (s *StateStorage) GetStateOrDefault(startPageURL string) (*State, error) {
	var state State
	result := s.db.First(&state, "start_page = ?", startPageURL)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &State{
			StartPage: startPageURL,
			NextPage:  startPageURL,
		}, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get previous page: %w", result.Error)
	}

	return &state, nil
}

func (s *StateStorage) SaveState(state *State) error {
	return s.db.Save(state).Error
}
