package usecase

import (
	"context"
	"log"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
)

type credentialUseCase struct {
	repo repo.IUserRepo
}

func NewCredentialUseCase(repo repo.IUserRepo) ICredentialUseCase {
	return &credentialUseCase{repo}
}

func (uc *credentialUseCase) SignUp(ctx context.Context, user domain.User) (domain.User, error) {
	if err := user.ValidateWithContext(ctx); err != nil {
		log.Println("Error occurred when validating")
		return user, err
	}
	if user, err := uc.repo.CreateNewUser(ctx, user); err != nil {
		log.Println("Error occured when persisting data")
		return user, err
	}

	return user, nil
}
