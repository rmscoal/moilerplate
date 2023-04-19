package mapper

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
)

func MapUserDomainToPersistence(user domain.User) model.User {
	res := model.User{
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
