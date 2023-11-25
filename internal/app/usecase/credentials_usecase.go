package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/repo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/service"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/utils"
)

type credentialUseCase struct {
	repo    repo.ICredentialRepo
	service service.IDoorkeeperService
}

func NewCredentialUseCase(repo repo.ICredentialRepo, service service.IDoorkeeperService) ICredentialUseCase {
	return &credentialUseCase{repo: repo, service: service}
}

func (uc *credentialUseCase) SignUp(ctx context.Context, user domain.User) (domain.User, error) {
	// Validate user entity
	if err := user.Validate(); err != nil {
		return user, NewDomainError("User", err)
	}

	// Validate repository state
	if err := uc.repo.ValidateRepoState(ctx, user); err != nil {
		return user, NewConflictError("User", err)
	}

	mixture, err := uc.service.HashPassword(user.Credential.Password)
	if err != nil {
		return user, NewServiceError("User", err)
	}
	user.Credential.SetEncodedPasswordFromByte(mixture)

	user, err = uc.repo.CreateNewUser(ctx, user)
	if err != nil {
		return user, NewRepositoryError("User", err)
	}

	user, err = uc.prepareUserTokensGeneration(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (uc *credentialUseCase) Login(ctx context.Context, cred vo.UserCredential) (domain.User, error) {
	var user domain.User

	// Validate request
	if err := cred.Validate(); err != nil {
		return user, NewDomainError("Credentials", err)
	}

	// Retrieve the user by its username
	user, err := uc.repo.GetUserByUsername(ctx, cred.Username)
	if err != nil {
		return user, NewNotFoundError("Credentials", err)
	}

	// Decode user's mixture
	mixture, err := user.Credential.GetHashMixture()
	if err != nil {
		return user, NewDomainError("Credentials", utils.AddError(fmt.Errorf("unable to retrieve user's hash"), err))
	}

	// Compare password request with users hashed password
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	success, err := uc.service.CompareHashAndPassword(ctxWithTimeout, cred.Password, mixture)
	if err != nil || !success {
		return user, NewUnauthorizedError(utils.AddError(fmt.Errorf("the password does not match"), err))
	}

	user, err = uc.prepareUserTokensGeneration(ctx, user)
	if err != nil {
		return user, err
	}

	go uc.generateNewHashMixture(user.Id, cred.Password)

	return user, nil
}

func (uc *credentialUseCase) Authorize(ctx context.Context, token string) (domain.User, error) {
	id, err := uc.service.VerifyAndParseToken(ctx, token)
	if err != nil {
		return domain.User{}, NewUnauthorizedError(err)
	}

	return domain.User{Id: id}, nil
}

// Refresh validates refresh tokens and generates a new set of tokens.
// Below are the steps:
//  1. Verify and parse incoming refresh token to retrieve the jti (JWT ID).
//  2. Checks whether the jti is present in the repository.
//     a. Case not exists: throw unauthorized error
//  3. Validate the jti:
//     a. Case it fails (reuse of refresh token... meaning a stolen one), then invalidates
//     all the refresh tokens families.
//
// 4. Generate a new version record on the repo
// 5. Generate the tokens
// See: https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/
func (uc *credentialUseCase) Refresh(ctx context.Context, refreshToken string) (domain.User, error) {
	jti, err := uc.service.VerifyAndParseRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.User{}, NewUnauthorizedError(err)
	}

	user, err := uc.repo.GetUserByJti(ctx, jti)
	if err != nil {
		return domain.User{}, NewUnauthorizedErrorWithReport(err)
	}

	if err := uc.validateUserToken(ctx, user); err != nil {
		return domain.User{}, err // usecase error
	}

	user, err = uc.prepareUserTokensGeneration(ctx, user)
	if err != nil {
		return user, err // usecase error
	}

	return user, nil
}

func (uc *credentialUseCase) prepareUserTokensGeneration(ctx context.Context, user domain.User) (domain.User, error) {
	token, err := uc.repo.SetNewUserToken(ctx, user)
	if err != nil {
		return domain.User{}, NewRepositoryError("Credentials", err)
	}

	user.Credential.Tokens = token
	token, err = uc.service.GenerateUserTokens(user)
	if err != nil {
		if rErr := uc.repo.UndoSetUserToken(ctx, user.Credential.Tokens.TokenID); rErr != nil {
			err = utils.AddError(err, rErr)
		}
		return domain.User{}, NewServiceError("Credentials", err)
	}

	user.Credential.Tokens = token
	return user, nil
}

func (uc *credentialUseCase) validateUserToken(ctx context.Context, user domain.User) error {
	version, err := uc.repo.GetLatestUserTokenVersion(ctx, user)
	if err != nil {
		return NewRepositoryError("Credentials", err)
	}

	if user.Credential.Tokens.Issued || user.Credential.Tokens.Version < version {
		err := fmt.Errorf("jti was issued before")
		rErr := uc.lockDownUser(ctx, user)
		if rErr != nil {
			err = utils.AddError(err, rErr)
		}
		return NewUnauthorizedError(err)
	}

	return nil
}

func (uc *credentialUseCase) lockDownUser(ctx context.Context, user domain.User) error {
	if err := uc.repo.DeleteUserTokenFamily(ctx, user); err != nil {
		return NewRepositoryError("Credentials", err)
	}
	return nil
}

// generateNewHashMixture method creates a new hash
// mixture from user's password. Set this up with
// goroutine after signin in.
func (uc *credentialUseCase) generateNewHashMixture(id, password string) {
	user := domain.User{Id: id, Credential: vo.UserCredential{Password: password}}

	mixture, err := uc.service.HashPassword(user.Credential.Password)
	if err != nil {
		log.Printf("failed to generate new hash password for user id: %s", id)
	}
	user.Credential.SetEncodedPasswordFromByte(mixture)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if err := uc.repo.RotateUserHashPassword(ctx, user); err != nil {
		log.Printf("failed to save user hash rotation for user id: %s", id)
	}
}
