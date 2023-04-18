package composer

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
)

type IUseCaseComposer interface {
	CredentialUseCase() usecase.ICredentialUseCase
}

type useCaseComposer struct {
	repo IRepoComposer
}

func NewUseCaseComposer(repo IRepoComposer) IUseCaseComposer {
	return &useCaseComposer{repo}
}

func (c *useCaseComposer) CredentialUseCase() usecase.ICredentialUseCase {
	return usecase.NewCredentialUseCase(c.repo.CredentialRepo())
}
