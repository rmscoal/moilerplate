package composer

import (
	"log"

	impl "github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/postgres"
)

type IRepoComposer interface {
	CredentialRepo() repo.ICredentialRepo
	UserProfileRepo() repo.IUserProfileRepo
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

	comp.Migrate()

	switch comp.env {
	case "DEVELOPMENT":
		comp.setToDebug()
	}

	return comp
}

// -------------- DI --------------
func (c *repoComposer) CredentialRepo() repo.ICredentialRepo {
	return impl.NewCredentialRepo()
}

func (c *repoComposer) UserProfileRepo() repo.IUserProfileRepo {
	return impl.NewUserProfileRepo()
}

// -------------- Setups --------------
func (c *repoComposer) setToDebug() {
	c.db.ORM = c.db.ORM.Debug()
}

func (c *repoComposer) Migrate() {
	if err := c.db.ORM.AutoMigrate(
		model.GetAllRelationalModels()...,
	); err != nil {
		log.Fatalf("FATAL - Unable to automigrate models: %s", err)
	}

	impl.InitBaseRepo(c.db.ORM)
}
