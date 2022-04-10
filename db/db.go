package db

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"corona-crawler/articles"
	"corona-crawler/history"
)

func New(postgresDSN string) (*gorm.DB, error) {
	// Replace default logger
	zLog := log.Logger.With().Caller().Logger()

	var gormDB *gorm.DB
	retries := 5
	for range time.Tick(time.Second) {
		db, err := gorm.Open(
			postgres.New(postgres.Config{DSN: postgresDSN}),
			&gorm.Config{
				Logger: logger.New(
					&zLog,
					logger.Config{
						SlowThreshold:             time.Second,
						LogLevel:                  logger.Warn,
						IgnoreRecordNotFoundError: true,
						Colorful:                  false,
					},
				),
			},
		)
		if err == nil {
			gormDB = db
			break
		}
		if retries < 0 {
			return nil, fmt.Errorf("failed to connect to database: no more retries: %w", err)
		}
		retries--
	}

	// TODO: Replace with proper migration (but should be enough)
	err := gormDB.AutoMigrate(&articles.ShortModel{}, &crawler.State{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return gormDB, err
}
