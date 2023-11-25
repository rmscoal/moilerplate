package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/mapper"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"gorm.io/gorm"
)

type userProfileRepo struct {
	db *gorm.DB
}

func NewUserProfileRepo(db *gorm.DB) *userProfileRepo {
	return &userProfileRepo{db}
}

func (repo *userProfileRepo) SaveUserEmails(ctx context.Context, user domain.User) (domain.User, error) {
	userModel := mapper.MapUserDomainToPersistence(user)
	tx := repo.db.WithContext(ctx).Begin()

	if err := tx.Unscoped().
		Delete(&model.UserEmail{}, "user_id = ?", user.Id).
		Error; err != nil {
		tx.Rollback()
		return user, translateGORMError(err)
	}

	if err := tx.Model(&userModel).
		Select("UserEmails").
		Save(&userModel).Error; err != nil {
		tx.Rollback()
		return user, translateGORMError(err)
	}
	tx.Commit()
	return mapper.MapUserModelToDomain(userModel), nil
}

func (repo *userProfileRepo) GetUserProfile(ctx context.Context, id string) (domain.User, error) {
	var userModel model.User
	if err := repo.db.WithContext(ctx).
		Model(&userModel).
		Preload("UserEmails").
		Preload("UserCredential").
		First(&userModel, "id = ?", id).
		Error; err != nil {
		return domain.User{}, translateGORMError(err)
	}

	return mapper.MapUserModelToDomain(userModel), nil
}
