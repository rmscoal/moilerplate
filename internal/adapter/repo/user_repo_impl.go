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

func (repo *userRepo) GetUserByCredentials(ctx context.Context, cred vo.UserCredential) (domain.User, error) {
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
