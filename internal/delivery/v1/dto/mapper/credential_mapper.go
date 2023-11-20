package mapper

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/dto"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

// Add type for namespace
type credentialMapper int

// Credential mapper namespace
var Credential credentialMapper

func (credentialMapper) SignUpRequestToUserDomain(obj dto.SignUpRequest) domain.User {
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

func (credentialMapper) LoginRequestToUserCredential(obj dto.LoginRequest) vo.UserCredential {
	return vo.UserCredential{
		Username: obj.Username,
		Password: obj.Password,
	}
}

func (credentialMapper) UserDomainToTokenResponse(user domain.User) dto.TokenResponse {
	return dto.TokenResponse{
		AccessToken:  user.Credential.Tokens.AccesssToken,
		RefreshToken: user.Credential.Tokens.RefreshToken,
		Username:     user.Credential.Username,
	}
}
