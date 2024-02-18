package usecase

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/app/repo"
	"github.com/rmscoal/moilerplate/internal/app/service"
	"github.com/rmscoal/moilerplate/internal/domain"
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

	user, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return user, NewRepositoryError("User", err)
	}

	return user, nil
}

func (uc *credentialUseCase) Login(ctx context.Context, user domain.User) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (uc *credentialUseCase) Authenticate(ctx context.Context, token string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (uc *credentialUseCase) Refresh(ctx context.Context, token string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}
