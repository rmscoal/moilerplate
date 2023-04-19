package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/mapper"
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

	// Probably should just attach the new id.
	user = mapper.MapUserModelToDomain(model)

	return user, nil
}
