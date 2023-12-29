package repo

import (
	"context"
	"fmt"

	"github.com/rmscoal/moilerplate/internal/adapter/repo/mapper"
	"github.com/rmscoal/moilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/moilerplate/internal/domain"
)

type userProfileRepo struct {
	*baseRepo
}

func NewUserProfileRepo() *userProfileRepo {
	return &userProfileRepo{baseRepo: gormRepo}
}

func (repo *userProfileRepo) SaveUserEmails(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, span := repo.tracer.Start(ctx, "(*userProfileRepo).SaveUserEmails")
	defer span.End()

	userModel := mapper.MapUserDomainToPersistence(user)
	tx := repo.db.WithContext(ctx).Begin()

	if err := tx.Unscoped().
		Delete(&model.UserEmail{}, "user_id = ?", user.Id).
		Error; err != nil {
		tx.Rollback()
		return user, fmt.Errorf("unable to delete all user emails")
	}

	if err := tx.Model(&userModel).
		Select("UserEmails").
		Save(&userModel).Error; err != nil {
		tx.Rollback()
		return user, repo.DetectConstraintError(err)
	}
	tx.Commit()
	return mapper.MapUserModelToDomain(userModel), nil
}

func (repo *userProfileRepo) GetUserProfile(ctx context.Context, id string) (domain.User, error) {
	ctx, span := repo.tracer.Start(ctx, "(*userProfileRepo).GetUserProfile")
	defer span.End()

	var userModel model.User
	if err := repo.db.WithContext(ctx).
		Model(&userModel).
		Preload("UserEmails").
		Preload("UserCredential").
		First(&userModel, "id = ?", id).
		Error; err != nil {
		return domain.User{}, repo.DetectNotFoundError(err)
	}

	return mapper.MapUserModelToDomain(userModel), nil
}
