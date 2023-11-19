package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/config"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/composer"
	v1 "github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/doorkeeper"
	httpserver "github.com/rmscoal/go-restful-monolith-boilerplate/pkg/http"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/logger"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/postgres"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/rater"
)

func Run(cfg *config.Config) {
	// Postgres .-.
	pg := postgres.GetPostgres(
		cfg.Db.URL,
		postgres.MaxPoolSize(cfg.Db.MaxPoolSize()),
		postgres.MaxOpenCoon(cfg.Db.MaxOpenConn()),
	)

	// Logger .-.
	logger := logger.NewAppLogger(cfg.App.LogPath())

	dk := doorkeeper.GetDoorkeeper(
		doorkeeper.RegisterHasherFunc(cfg.Doorkeeper.HashMethod()),
		doorkeeper.RegisterSignMethod(cfg.Doorkeeper.SigningMethod(), cfg.Doorkeeper.SignSize()),
		doorkeeper.RegisterIssuer(cfg.Doorkeeper.Issuer()),
		doorkeeper.RegisterAccessDuration(cfg.Doorkeeper.AccessTokenDuration()),
		doorkeeper.RegisterRefreshDuration(cfg.Doorkeeper.RefreshTokenDuration()),
		doorkeeper.RegisterCertPath(cfg.Doorkeeper.CertPath()),
		doorkeeper.RegisterSecretKey(cfg.Doorkeeper.SecretKey()),
	)

	// Rate Limitter .-.
	rt := rater.GetRater(context.Background(),
		rater.RegisterRateLimitForEachClient(cfg.App.RaterLimit()),
		rater.RegisterBurstLimitForEachClient(cfg.App.BurstLimit()),
		rater.RegisterEvaluationInterval(cfg.App.RaterEvaluationInterval()),
		rater.RegisterDeletionTime(cfg.App.RaterDeletionTime()),
	)

	// Composers .-.
	serviceComposer := composer.NewServiceComposer(dk, rt)
	repoComposer := composer.NewRepoComposer(pg, cfg.App.Environment())
	usecaseComposer := composer.NewUseCaseComposer(repoComposer, serviceComposer)

	// Http
	deliveree := gin.Default()
	v1.NewRouter(deliveree, logger, usecaseComposer)
	httpserver.NewServer(deliveree,
		httpserver.RegisterHostAndPort(cfg.Server.Host, cfg.Server.Port),
	)
}
