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
	db *postgres.Postgres
}

func NewRepoComposer(db *postgres.Postgres) IRepoComposer {
	comp := new(repoComposer)
	comp.db = db

	comp.Migrate()

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
func (c *repoComposer) Migrate() {
	if err := c.db.ORM.AutoMigrate(
		model.GetAllRelationalModels()...,
	); err != nil {
		log.Fatalf("FATAL - Unable to automigrate models: %s", err)
	}

	impl.InitBaseRepo(c.db.ORM)
}
