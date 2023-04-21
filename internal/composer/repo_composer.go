package composer

import (
	impl "github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/postgres"
)

type IRepoComposer interface {
	CredentialRepo() repo.ICredentialRepo
	Migrate()
}

type repoComposer struct {
	db  *postgres.Postgres
	env string
}

func NewRepoComposer(db *postgres.Postgres, env string) IRepoComposer {
	comp := new(repoComposer)
	comp.env = env
	comp.db = db

	switch comp.env {
	case "DEVELOPMENT":
		comp.setToDebug()
	case "MIGRATION":
		comp.setToDebug()
		comp.Migrate()
	}

	return comp
}

func (c *repoComposer) setToDebug() {
	c.db.ORM = c.db.ORM.Debug()
}

func (c *repoComposer) CredentialRepo() repo.ICredentialRepo {
	return impl.NewCredentialRepo(c.db.ORM)
}

func (c *repoComposer) Migrate() {
	c.db.ORM.AutoMigrate(
		model.GetAllRelationalModels()...,
	)
}
