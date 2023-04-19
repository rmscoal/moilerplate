package composer

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
)

type IUseCaseComposer interface {
	CredentialUseCase() usecase.ICredentialUseCase
}

type useCaseComposer struct {
	repo    IRepoComposer
	service IServiceComposer
}

func NewUseCaseComposer(repo IRepoComposer, service IServiceComposer) IUseCaseComposer {
	return &useCaseComposer{repo: repo, service: service}
}

func (c *useCaseComposer) CredentialUseCase() usecase.ICredentialUseCase {
	return usecase.NewCredentialUseCase(c.repo.CredentialRepo(), c.service.DoorkeeperService())
}
