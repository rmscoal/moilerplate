package composer

import (
	"log"

	impl "github.com/rmscoal/moilerplate/internal/adapter/repo"
	"github.com/rmscoal/moilerplate/internal/app/repo"
	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/pkg/postgres"
	"gorm.io/plugin/opentelemetry/tracing"
)

type IRepoComposer interface {
	CredentialRepo() repo.ICredentialRepo
	Migrate()
}

type repoComposer struct {
	db *postgres.Postgres
}

func NewRepoComposer(db *postgres.Postgres) IRepoComposer {
	comp := new(repoComposer)
	comp.db = db

	comp.Migrate()
	// Use tracing after migrating
	comp.db.ORM.Use(tracing.NewPlugin(tracing.WithoutQueryVariables(), tracing.WithoutMetrics()))

	return comp
}

// -------------- DI --------------

func (c *repoComposer) CredentialRepo() repo.ICredentialRepo {
	return impl.NewCredentialRepo()
}

// -------------- Setups --------------

func (c *repoComposer) Migrate() {
	if err := c.db.ORM.AutoMigrate(
		domain.User{},
	); err != nil {
		log.Fatalf("FATAL - Unable to automigrate models: %s", err)
	}

	impl.InitBaseRepo(c.db.ORM, true)
}
