package mapper

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

func MapModifyEmailRequestToUserDomain(id string, obj dto.ModifyEmailRequest) domain.User {
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

func MapUserDomainToModifyEmailResponse(user domain.User) dto.ModifyEmailResponse {
	obj := dto.ModifyEmailResponse{
		UserId: user.Id,
	}

	for _, email := range user.Emails {
		obj.Emails = append(obj.Emails, dto.ModifyEmailDetailRequest{
			Email:     email.Email,
			IsPrimary: email.IsPrimary,
		})
	}

	return obj
}
