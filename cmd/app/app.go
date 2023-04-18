package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/config"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/composer"
	v1 "github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1"
	httpserver "github.com/rmscoal/go-restful-monolith-boilerplate/pkg/http"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/logger"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/postgres"
)

func Run(cfg *config.Config) {
	// Postgres .-.
	pg := postgres.GetPostgres(
		cfg.Db.URL,
		postgres.MaxPoolSize(cfg.Db.MaxPoolSize()),
		postgres.MaxOpenCoon(cfg.Db.MaxOpenConn()),
	)

	// Logger .-.
	logger := logger.NewAppLogger(cfg.App.LogPath)

	// Composers .-.
	repoComposer := composer.NewRepoComposer(pg, cfg.App.Environment)
	usecaseComposer := composer.NewUseCaseComposer(repoComposer)

	// Http
	deliveree := gin.Default()
	v1.NewRouter(deliveree, logger, usecaseComposer)
	httpserver.NewServer(deliveree)
}
