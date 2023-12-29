package usecase

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/app/repo"
	"github.com/rmscoal/moilerplate/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type userProfileUseCase struct {
	repo   repo.IUserProfileRepo
	tracer trace.Tracer
}

func NewUserProfileUseCase(repo repo.IUserProfileRepo) IUserProfileUseCase {
	return &userProfileUseCase{repo: repo, tracer: otel.Tracer("profile_usecase")}
}

func (uc *userProfileUseCase) RetrieveProfile(ctx context.Context, userID string) (domain.User, error) {
	ctx, span := uc.tracer.Start(ctx, "(*userProfileUseCase).RetrieveProfile")
	defer span.End()

	user, err := uc.repo.GetUserProfile(ctx, userID)
	if err != nil {
		return user, NewNotFoundError("User", err)
	}

	return user, nil
}

func (uc *userProfileUseCase) ModifyEmailAddress(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, span := uc.tracer.Start(ctx, "(*userProfileUseCase).ModifyEmailAddress")
	defer span.End()

	if err := user.ValidateEmails(); err != nil {
		return user, NewDomainError("User", err)
	}

	// Persist to repo
	user, err := uc.repo.SaveUserEmails(ctx, user)
	if err != nil {
		return user, NewConflictError("User", err)
	}

	return user, nil
}
