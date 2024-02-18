package mapper

import (
	"github.com/rmscoal/moilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/moilerplate/internal/domain"
)

type credentialMapper int

var Credential credentialMapper

func (credentialMapper) SignupRequestToUserDomain(req dto.SignUpRequest) domain.User {
	return domain.User{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}
}
