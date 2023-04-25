package mapper

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

func MapUserDomainToNewAuthCredModel(user domain.User) model.AuthorizationCredential {
	res := model.AuthorizationCredential{
		BaseModelId: model.BaseModelId{},
		UserId:      user.Id,
		Version:     user.Credential.Tokens.Version + 1,
	}

	if user.Credential.Tokens.TokenID == "" {
		return res
	}

	res.ParentId = &user.Credential.Tokens.TokenID
	return res
}

func MapAuthCredToUserTokenVO(authCred model.AuthorizationCredential) vo.UserToken {
	return vo.UserToken{
		TokenID:  authCred.Id,
		Version:  authCred.Version,
		Issued:   authCred.Issued,
		IssuedAt: authCred.IssuedAt,
	}
}

func MapAuthCredModelToUserDomain(authCred model.AuthorizationCredential) domain.User {
	tokens := MapAuthCredToUserTokenVO(authCred)
	res := domain.User{
		Id:        authCred.UserId,
		FirstName: authCred.User.FirstName,
		LastName:  authCred.User.LastName,
		Credential: vo.UserCredential{
			Username: authCred.User.UserCredential.Username,
		},
	}

	res.Credential.Tokens = tokens
	return res
}
