package repo

import (
	"context"
	"fmt"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/mapper"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"gorm.io/gorm"
)

type credentialRepo struct {
	db *gorm.DB
}

func NewCredentialRepo(db *gorm.DB) *credentialRepo {
	return &credentialRepo{db}
}

func (repo *credentialRepo) CreateNewUser(ctx context.Context, user domain.User) (domain.User, error) {
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

func (repo *credentialRepo) GetUserByCredentials(ctx context.Context, cred vo.UserCredential) (domain.User, error) {
	var userModel model.User

	if err := repo.db.
		WithContext(ctx).
		Model(&userModel).
		InnerJoins("UserCredential", repo.db.Where(&model.UserCredential{Username: cred.Username, Password: cred.Password})).
		First(&userModel).
		Error; err != nil {
		return domain.User{}, fmt.Errorf("user not found with given username and password")
	}
	// Alternative:
	// if err := repo.db.
	// 	WithContext(ctx).
	// 	Model(&userModel).
	// 	Joins(`INNER JOIN user_credentials ON user_credentials.user_id = users.id`).
	// 	Where("user_credentials.username = ?", cred.Username).
	// 	Where("user_credentials.password = ?", cred.Password).
	// 	First(&userModel).
	// 	Error; err != nil {
	// 	return domain.User{}, fmt.Errorf("user not found with given username and password")
	// }

	return mapper.MapUserModelToDomain(userModel), nil
}

func (repo *credentialRepo) GetUserByJti(ctx context.Context, jti string) (domain.User, error) {
	var userModel model.User
	authCred := model.AuthorizationCredential{BaseModelId: model.BaseModelId{Id: jti}}

	if err := repo.db.WithContext(ctx).
		Where(&authCred).
		Association("User").
		Find(&userModel); err != nil {
		return domain.User{}, fmt.Errorf("unable to get user from jti")
	}

	if err := repo.db.WithContext(ctx).Model(&authCred).Update("issued", true).Error; err != nil {
		return domain.User{}, fmt.Errorf("unable to issue jti")
	}

	userModel.AuthorizationCredentials = append(userModel.AuthorizationCredentials, authCred)
	return mapper.MapUserModelToDomain(userModel), nil
}

func (repo *credentialRepo) SetNewUserToken(ctx context.Context, user domain.User) (vo.UserToken, error) {
	authCred := mapper.MapUserDomainToNewAuthCredModel(user)

	fmt.Printf("%+v\n", authCred)

	if err := repo.db.WithContext(ctx).Create(&authCred).Error; err != nil {
		return vo.UserToken{}, fmt.Errorf("unable to set a new user token: %s", err)
	}

	return mapper.MapAuthCredToUserTokenVO(authCred), nil
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

/*
*************************************************
REPO VALIDATIONS IMPLEMENTATIONS
*************************************************
*/
func (repo *credentialRepo) ValidateRepoState(ctx context.Context, user domain.User) error {
	var err error
	if repo.UsernameExists(ctx, user.Id, user.Credential.Username) {
		err = AddError(err, fmt.Errorf("username has been taken"))
	}
	if repo.EmailExists(ctx, user.Id, user.Emails[0].Email) {
		err = AddError(err, fmt.Errorf("email has been taken"))
	}
	return err
}

func (repo *credentialRepo) UsernameExists(ctx context.Context, id string, username string) bool {
	var userId string
	repo.db.WithContext(ctx).
		Model(&model.UserCredential{}).
		Select("user_id").
		Where("username = ?", username).
		Scan(&userId)
	return userId != id
}

func (repo *credentialRepo) EmailExists(ctx context.Context, id string, email string) bool {
	var userId string
	repo.db.WithContext(ctx).
		Model(&model.UserEmail{}).
		Select("user_id").
		Where("email = ?", email).
		Scan(&userId)
	return userId != id
}
