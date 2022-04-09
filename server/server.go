package server

import (
	"corona-crawler/crawler"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"time"
)

const postgresDSN = "postgres://corona-user:corona-password@postgres/corona-crawler"

func Run() error {
	db, err := crawler.NewDB(postgresDSN)
	if err != nil {
		return err
	}

	// TODO: Should replace logger
	e := echo.New()
	e.HideBanner = true

	e.GET("/", func(c echo.Context) error {
		yearAgo := time.Now().AddDate(-1, 0, 0)

		var count int64
		err := db.Model(&crawler.ArticleModel{}).Where("published > ?", yearAgo).Count(&count).Error
		if err != nil {
			log.Error().Err(err).Msg("failed to count articles")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, struct {
			Count int64 `json:"count"`
		}{
			Count: count,
		})
	})

	return e.Start(":8080")
}
