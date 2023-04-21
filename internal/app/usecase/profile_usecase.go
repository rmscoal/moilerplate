package usecase

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
)

type userProfileUseCase struct {
	repo repo.IUserProfileRepo
}

func NewUserProfileUseCase(repo repo.IUserProfileRepo) IUserProfileUseCase {
	return &userProfileUseCase{repo}
}

func (uc *userProfileUseCase) ModifyEmailAddress(ctx context.Context, userReq domain.User) (domain.User, error) {
	// TODO:
	// Should probably load first the user from repo...
	// Checks which is the new and old emails...
	// Verify email by sending OTP code...

	user, err := uc.repo.GetUserProfile(ctx, userReq.Id)
	if err != nil {
		return user, NewNotFoundError("User", err)
	}

	user.ReplaceEmails(userReq.Emails)

	// Validate each emails matches domain rules
	if err := user.ValidateEmailsWithContext(ctx); err != nil {
		return user, NewDomainError("User", err)
	}

	// Verify that there exists only one primary email
	if err := user.VerifyOnePrimaryEmail(); err != nil {
		return user, NewDomainError("User", err)
	}

	// Persist to repo
	user, err = uc.repo.SaveUserEmails(ctx, user)
	if err != nil {
		return user, NewConflictError("User", err)
	}

	return user, nil
}
