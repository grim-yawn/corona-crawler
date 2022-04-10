package articles

import "time"

type ArticleID int

type ShortModel struct {
	ID        ArticleID `gorm:"primaryKey;autoIncrement:false"`
	Published time.Time
}

func (m ShortModel) TableName() string {
	return "articles"
}

type Article struct {
	ID ArticleID `json:"id"`

	Content struct {
		Title     string `json:"title"`
		TitleHead string `json:"titleHead"`
		Lead      string `json:"lead"`

		Published time.Time `json:"published"`
	} `json:"content"`
}

func (a Article) ToModel() ShortModel {
	return ShortModel{
		ID:        a.ID,
		Published: a.Content.Published,
	}
}

func ToModels(articles []Article) []ShortModel {
	models := make([]ShortModel, 0, len(articles))
	for _, article := range articles {
		models = append(models, article.ToModel())
	}
	return models
}
