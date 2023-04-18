package mapper

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

func SignUpRequestToUserDomain(obj dto.SignUpRequest) domain.User {
	return domain.User{
		FirstName: obj.FirstName,
		LastName:  obj.LastName,
		Emails: []vo.UserEmail{
			{
				Email:     obj.Email,
				IsPrimary: true,
			},
		},
		PhoneNumber: obj.PhoneNumber,
		Credential: vo.UserCredential{
			Username: obj.Username,
			Password: obj.Password,
		},
	}
}
