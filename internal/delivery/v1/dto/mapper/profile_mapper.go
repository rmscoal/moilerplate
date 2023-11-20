package mapper

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

// Profile mapper namespace
type profileMapper int

// Namespace to call
var Profile profileMapper

func (profileMapper) MapUserDomainToFullProfileResponse(user domain.User) dto.FullProfileResponse {
	profile := dto.FullProfileResponse{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Credential.Username,
		PhoneNumber: user.PhoneNumber,
		Emails:      []dto.EmailDetail{},
	}

	for _, email := range user.Emails {
		profile.Emails = append(profile.Emails, dto.EmailDetail{
			Email:     email.Email,
			IsPrimary: email.IsPrimary,
		})
	}

	return profile
}

func (profileMapper) MapModifyEmailRequestToUserDomain(id string, obj dto.ModifyEmailRequest) domain.User {
	user := domain.User{
		Id: id,
	}
	for _, email := range obj.Emails {
		user.Emails = append(user.Emails, vo.UserEmail{
			Email:     email.Email,
			IsPrimary: email.IsPrimary,
		})
	}
	return user
}

func (profileMapper) MapUserDomainToModifyEmailResponse(user domain.User) dto.ModifyEmailResponse {
	obj := dto.ModifyEmailResponse{
		UserId: user.Id,
	}

	for _, email := range user.Emails {
		obj.Emails = append(obj.Emails, dto.EmailDetail{
			Email:     email.Email,
			IsPrimary: email.IsPrimary,
		})
	}

	return obj
}
