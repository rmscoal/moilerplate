package usecase

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/app/repo"
	"github.com/rmscoal/moilerplate/internal/domain"
)

type userProfileUseCase struct {
	repo repo.IUserProfileRepo
}

func NewUserProfileUseCase(repo repo.IUserProfileRepo) IUserProfileUseCase {
	return &userProfileUseCase{repo}
}

func (uc *userProfileUseCase) RetrieveProfile(ctx context.Context, userID string) (domain.User, error) {
	user, err := uc.repo.GetUserProfile(ctx, userID)
	if err != nil {
		return user, NewNotFoundError("User", err)
	}

	return user, nil
}

func (uc *userProfileUseCase) ModifyEmailAddress(ctx context.Context, user domain.User) (domain.User, error) {
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
