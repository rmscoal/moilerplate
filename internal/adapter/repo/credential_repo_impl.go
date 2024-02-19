package repo

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
	"go.opentelemetry.io/otel/codes"
)

type credentialRepo struct {
	*baseRepo
}

func NewCredentialRepo() *credentialRepo {
	return &credentialRepo{baseRepo: gormRepo}
}

func (repo *credentialRepo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, span := repo.tracer.Start(ctx, "repo.CreateUser")
	defer span.End()

	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		span.SetStatus(codes.Error, "unable to create user")
		span.RecordError(err)
		return user, repo.DetectConstraintError(err)
	}

	return user, nil
}

func (repo *credentialRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	ctx, span := repo.tracer.Start(ctx, "repo.GetUserByUsername")
	defer span.End()

	var user domain.User
	if err := repo.db.WithContext(ctx).Take(&user, "username = ?", username).Error; err != nil {
		return user, repo.DetectNotFoundError(err)
	}

	return user, nil
}
