package usecase

import (
	"context"
	"errors"

	"github.com/rmscoal/moilerplate/internal/app/repo"
	"github.com/rmscoal/moilerplate/internal/app/service"
	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type credentialUseCase struct {
	repo    repo.ICredentialRepo
	service service.IDoorkeeperService
	tracer  trace.Tracer
}

func NewCredentialUseCase(repo repo.ICredentialRepo, service service.IDoorkeeperService) ICredentialUseCase {
	return &credentialUseCase{repo: repo, service: service, tracer: otel.Tracer("credential_usecase")}
}

func (uc *credentialUseCase) SignUp(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.SignUp")
	defer span.End()

	if err := user.Validate(); err != nil {
		return user, NewDomainError("User", err)
	}

	user.Password = uc.service.HashAndEncodeStringWithSalt(ctx, user.Password, user.Username)

	user, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return user, NewRepositoryError("User", err)
	}

	return user, nil
}

func (uc *credentialUseCase) Login(ctx context.Context, cred domain.User) (token vo.Token, err error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.Login")
	defer span.End()

	if err := cred.ValidateCredential(); err != nil {
		return token, err
	}

	user, err := uc.repo.GetUserByUsername(ctx, cred.Username)
	if err != nil {
		return token, NewUnauthorizedError(err)
	}

	match, err := uc.service.ComparePasswords(ctx, user.Password, cred.Password, user.Username)
	if err != nil {
		return token, NewServiceError("Doorkeeper", err)
	} else if !match {
		return token, NewUnauthorizedError(errors.New("password mismatched"))
	}

	token, err = uc.service.GenerateTokens(ctx, user.ID, nil)
	if err != nil {
		return token, NewServiceError("Doorkeeper", err)
	}

	return
}

func (uc *credentialUseCase) Authenticate(ctx context.Context, token string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (uc *credentialUseCase) Refresh(ctx context.Context, token string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}
