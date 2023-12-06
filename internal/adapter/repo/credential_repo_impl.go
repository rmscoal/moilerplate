package repo

import (
	"context"
	"fmt"

	"github.com/rmscoal/moilerplate/internal/adapter/repo/mapper"
	"github.com/rmscoal/moilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
	"gorm.io/gorm"
)

type credentialRepo struct {
	*baseRepo
}

func NewCredentialRepo() *credentialRepo {
	return &credentialRepo{baseRepo: gormRepo}
}

func (repo *credentialRepo) CreateNewUser(ctx context.Context, user domain.User) (domain.User, error) {
	model := mapper.MapUserDomainToPersistence(user)
	if err := repo.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		WithContext(ctx).
		Model(&model).
		Create(&model).Error; err != nil {
		return user, repo.TranslateError(err)
	}
	user = mapper.MapUserModelToDomain(model)
	return user, nil
}

func (repo *credentialRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var userModel model.User

	if err := repo.db.
		WithContext(ctx).
		Model(&userModel).
		InnerJoins("UserCredential", repo.db.Where(&model.UserCredential{Username: username})).
		First(&userModel).
		Error; err != nil {
		return domain.User{}, repo.TranslateError(err)
	}

	return mapper.MapUserModelToDomain(userModel), nil
}

func (repo *credentialRepo) GetUserByJti(ctx context.Context, jti string) (domain.User, error) {
	var authCred model.AuthorizationCredential

	if err := repo.db.WithContext(ctx).
		Preload("User").
		Preload("User.UserCredential").
		First(&authCred, "id = ?", jti).Error; err != nil {
		return domain.User{}, fmt.Errorf("unable to get user from jti: %s", err)
	}

	return mapper.MapAuthCredModelToUserDomain(authCred), nil
}

func (repo *credentialRepo) SetNewUserToken(ctx context.Context, user domain.User) (vo.UserToken, error) {
	authCred := mapper.MapUserDomainToNewAuthCredModel(user)

	if err := repo.IssueParentToken(ctx, authCred); err != nil {
		return vo.UserToken{}, err
	}

	if err := repo.db.WithContext(ctx).Create(&authCred).Error; err != nil {
		return vo.UserToken{}, fmt.Errorf("unable to set a new user token: %s", err)
	}

	return mapper.MapAuthCredToUserTokenVO(authCred), nil
}

func (repo *credentialRepo) IssueParentToken(ctx context.Context, authCred model.AuthorizationCredential) error {
	if authCred.ParentId != nil {
		if err := repo.db.WithContext(ctx).
			Model(&model.AuthorizationCredential{}).
			Where("id = ?", authCred.ParentId).
			Update("issued", true).Error; err != nil {
			return fmt.Errorf("unable to issue jti")
		}
	}

	return nil
}

func (repo *credentialRepo) UndoSetUserToken(ctx context.Context, jti string) error {
	if err := repo.db.WithContext(ctx).Delete(&model.AuthorizationCredential{}, "id = ?", jti).Error; err != nil {
		return fmt.Errorf("unable to undo creation of user token: %s", err)
	}
	return nil
}

func (repo *credentialRepo) GetLatestUserTokenVersion(ctx context.Context, user domain.User) (int, error) {
	var count int64
	if err := repo.db.WithContext(ctx).
		Model(&model.AuthorizationCredential{}).
		Where("user_id = ?", user.Id).
		Count(&count).
		Error; err != nil {
		return int(count), fmt.Errorf("unable to get latest version of token family: %s", err)
	}
	return int(count), nil
}

func (repo *credentialRepo) DeleteUserTokenFamily(ctx context.Context, user domain.User) error {
	if err := repo.db.WithContext(ctx).
		Delete(&model.AuthorizationCredential{}, "user_id = ?", user.Id).
		Error; err != nil {
		return fmt.Errorf("unable to invalidate user token family")
	}

	return nil
}

func (repo *credentialRepo) RotateUserHashPassword(ctx context.Context, user domain.User) error {
	tx := repo.db.WithContext(ctx).Begin()

	if err := tx.Model(&model.UserCredential{}).
		Where(&model.UserCredential{UserId: user.Id}).
		Update("password", user.Credential.Password).
		Error; err != nil {
		tx.Rollback()
		return repo.TranslateError(err)
	}

	tx.Commit()
	return nil
}

/*
*************************************************
REPO VALIDATIONS IMPLEMENTATIONS
*************************************************
*/
// Add your db state validations here
