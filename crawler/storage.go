package crawler

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type ArticleModel struct {
	ID        ArticleID `gorm:"primaryKey;autoIncrement:false"`
	Published time.Time
}

func (m ArticleModel) TableName() string {
	return "articles"
}

// PageModel tracks next page for crawling to avoid redundant crawling
type PageModel struct {
	StartPage string `gorm:"primaryKey"`

	NextPage string
}

func (p PageModel) TableName() string {
	return "page"
}

func NewDB(postgresDSN string) (*gorm.DB, error) {
	// Replace default logger
	zLog := log.Logger.With().Caller().Logger()

	// TODO: Should be replaced with proper retry connection
	db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{
		Logger: logger.New(
			&zLog,
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// TODO: Replace with proper migration (but should be enough)
	err = db.AutoMigrate(&ArticleModel{}, &PageModel{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return db, nil
}
