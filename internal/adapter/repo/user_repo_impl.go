package repo

import (
	"context"
	"fmt"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/mapper"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db}
}

func (repo *userRepo) CreateNewUser(ctx context.Context, user domain.User) (domain.User, error) {
	model := mapper.MapUserDomainToPersistence(user)
	if err := repo.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		WithContext(ctx).
		Model(&model).
		Create(&model).Error; err != nil {
		return user, err
	}
	user = mapper.MapUserModelToDomain(model)
	return user, nil
}

/*
*************************************************
REPO VALIDATIONS IMPLEMENTATIONS
*************************************************
*/
func (repo *userRepo) ValidateRepoState(ctx context.Context, user domain.User) error {
	var err error
	if repo.UsernameExists(ctx, user.Id, user.Credential.Username) {
		err = AddError(err, fmt.Errorf("username has been taken"))
	}
	if repo.EmailExists(ctx, user.Id, user.Emails[0].Email) {
		err = AddError(err, fmt.Errorf("email has been taken"))
	}
	return err
}

func (repo *userRepo) UsernameExists(ctx context.Context, id string, username string) bool {
	var userId string
	repo.db.WithContext(ctx).
		Model(&model.UserCredential{}).
		Select("user_id").
		Where("username = ?", username).
		Scan(&userId)
	return userId != id
}

func (repo *userRepo) EmailExists(ctx context.Context, id string, email string) bool {
	var userId string
	repo.db.WithContext(ctx).
		Model(&model.UserEmail{}).
		Select("user_id").
		Where("email = ?", email).
		Scan(&userId)
	return userId != id
}
