package mapper

import (
	"github.com/rmscoal/moilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
)

func MapUserDomainToPersistence(user domain.User) model.User {
	res := model.User{
		BaseModelId: model.BaseModelId{Id: user.Id},
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		UserCredential: model.UserCredential{
			Username: user.Credential.Username,
			Password: user.Credential.Password,
		},
	}

	for _, e := range user.Emails {
		email := model.UserEmail{
			Email:     e.Email,
			IsPrimary: e.IsPrimary,
		}

		res.UserEmails = append(res.UserEmails, email)
	}

	return res
}

func MapUserModelToDomain(user model.User) domain.User {
	res := domain.User{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Credential: vo.UserCredential{
			Username: user.UserCredential.Username,
			Password: user.UserCredential.Password,
		},
	}

	for _, e := range user.UserEmails {
		email := vo.UserEmail{
			Email:     e.Email,
			IsPrimary: e.IsPrimary,
		}
		res.Emails = append(res.Emails, email)
	}

	for _, cred := range user.AuthorizationCredentials {
		res.Credential.Tokens.TokenID = cred.Id
	}

	return res
}
