package usecase

import (
	"context"

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

// Login handle user signin for first time user and generate pair of a jwts
func (uc *credentialUseCase) SignUp(ctx context.Context, user domain.User) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

// Login handle user login and generate pair of jwts
func (uc *credentialUseCase) Login(ctx context.Context, cred vo.UserCredential) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

// Authenticates authenticates user from the given jwt.
func (uc *credentialUseCase) Authenticate(ctx context.Context, token string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

// Refresh validates refresh tokens and generates a new set of tokens.
func (uc *credentialUseCase) Refresh(ctx context.Context, token string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}
