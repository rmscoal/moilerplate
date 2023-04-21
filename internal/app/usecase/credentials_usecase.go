package usecase

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/service"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type credentialUseCase struct {
	repo    repo.ICredentialRepo
	service service.IDoorkeeperService
}

func NewCredentialUseCase(repo repo.ICredentialRepo, service service.IDoorkeeperService) ICredentialUseCase {
	return &credentialUseCase{repo: repo, service: service}
}

func (uc *credentialUseCase) SignUp(ctx context.Context, user domain.User) (domain.User, error) {
	// Validate user entity
	if err := user.ValidateWithContext(ctx); err != nil {
		return user, NewDomainError("User", err)
	}

	// Validate repository state
	if err := uc.repo.ValidateRepoState(ctx, user); err != nil {
		return user, NewConflictError("User", err)
	}

	user.Credential.Password = uc.service.HashPassword(user.Credential.Password)
	user, err := uc.repo.CreateNewUser(ctx, user)
	if err != nil {
		return user, NewRepositoryError("User", err)
	}

	token, err := uc.service.GenerateToken(user)
	if err != nil {
		return user, err
	}
	user.Credential.Token = token

	return user, nil
}

func (uc *credentialUseCase) Login(ctx context.Context, cred vo.UserCredential) (domain.User, error) {
	var user domain.User
	if err := cred.Validate(); err != nil {
		return user, NewDomainError("Credentials", err)
	}

	cred.Password = uc.service.HashPassword(cred.Password)
	user, err := uc.repo.GetUserByCredentials(ctx, cred)
	if err != nil {
		return user, NewNotFoundError("Credentials", err)
	}

	token, err := uc.service.GenerateToken(user)
	if err != nil {
		return user, err
	}
	user.Credential.Token = token
	return user, nil
}

func (uc *credentialUseCase) Authorize(ctx context.Context, token string) (domain.User, error) {
	id, err := uc.service.VerifyAndParseToken(token)
	if err != nil {
		return domain.User{}, NewUnauthorizedError(err)
	}

	return domain.User{Id: id}, nil
}
