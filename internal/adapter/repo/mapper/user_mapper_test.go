package mapper

import (
	"testing"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/repo/model"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserMapperTestSuite struct {
	suite.Suite
	user domain.User
}

func (m *UserMapperTestSuite) SetupTest() {
	m.user = domain.User{
		Id:        "",
		FirstName: "User",
		LastName:  "Test",
		Emails: []vo.UserEmail{
			{
				Email:     "test1@example.com",
				IsPrimary: true,
			},
		},
		PhoneNumber: "081234 567899",
		Credential: vo.UserCredential{
			Username: "username",
			Password: "password",
		},
	}
}

func TestUserMapperSuite(t *testing.T) {
	suite.Run(t, new(UserMapperTestSuite))
}

func (m *UserMapperTestSuite) TestMapUserDomainToPersistence() {
	assert.Equal(m.T(), MapUserDomainToPersistence(m.user), model.User{
		BaseModelId: model.BaseModelId{Id: m.user.Id},
		FirstName:   m.user.FirstName,
		LastName:    m.user.LastName,
		PhoneNumber: m.user.PhoneNumber,
		UserCredential: model.UserCredential{
			UserId:   m.user.Id,
			Username: m.user.Credential.Username,
			Password: m.user.Credential.Password,
		},
		UserEmails: []model.UserEmail{
			{
				UserId:    m.user.Id,
				Email:     m.user.Emails[0].Email,
				IsPrimary: m.user.Emails[0].IsPrimary,
			},
		},
	})
}
