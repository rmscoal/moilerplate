package service

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/service"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/doorkeeper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DoorkeeperServiceImplTestSuite struct {
	suite.Suite
	dk      *doorkeeper.Doorkeeper
	service service.IDoorkeeperService
}

func (suite *DoorkeeperServiceImplTestSuite) SetupTest() {
	dk := doorkeeper.GetDoorkeeper(
		doorkeeper.RegisterHasherFunc("SHA384"),
		doorkeeper.RegisterSignMethod("RSA", "256"),
		doorkeeper.RegisterIssuer("TESTAPP"),
		doorkeeper.RegisterAccessDuration(time.Duration(5*time.Minute)),
		doorkeeper.RegisterRefreshDuration(20*time.Minute),
		doorkeeper.RegisterCertPath("../../../cert"),
	)

	suite.dk = dk
	suite.service = NewDoorkeeperService(suite.dk)
}

func TestDoorkeeperService(t *testing.T) {
	suite.Run(t, new(DoorkeeperServiceImplTestSuite))
}

func (suite *DoorkeeperServiceImplTestSuite) TestHashPassword() {
	suite.Run("Successful HashPassword", func() {
		// Generate Hash Mixture
		hash, err := suite.service.HashPassword("password")
		assert.Nil(suite.T(), err)
		assert.Greater(suite.T(), len(hash), int(MinSaltLength))
		assert.LessOrEqual(suite.T(), len(hash), int(MaxSaltLength)+suite.dk.GetHashKeyLen())

		// Base64 Encode
		encoded := base64.StdEncoding.EncodeToString(hash)
		assert.Equal(suite.T(), rune(encoded[len(encoded)-1]), rune('='))
	})

	suite.Run("Panic on Generating Random Range", func() {
		// It seems that HashPassword will only throw error iff
		// MinSaltLength > MaxSaltLength or rand.Reader fails on
		// the os side.
		MinSaltLength, MaxSaltLength = MaxSaltLength, MinSaltLength // Intentional produce errors

		var f assert.PanicTestFunc
		f = func() {
			suite.service.HashPassword("password")
		}
		assert.Panics(suite.T(), f)
	})
}

func (suite *DoorkeeperServiceImplTestSuite) TestCompareHashAndPassword() {
	// Makes sure that generating hash runs as expected
	MinSaltLength = 1 << 5
	MaxSaltLength = 1 << 6

	suite.Run("Successful CompareHashAndPassword", func() {
		hash, _ := suite.service.HashPassword("password")

		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()
		found, err := suite.service.CompareHashAndPassword(ctx, "password", hash)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), found, true)
	})

	suite.Run("Unsuccesful CompareHashAndPassword", func() {
		suite.Run("Invalid Hash Length", func() {
			hash := make([]byte, 0)

			ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
			defer cancel()
			found, err := suite.service.CompareHashAndPassword(ctx, "wrong_password", hash)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "invalid hash length")
			assert.ErrorIs(suite.T(), err, ErrInvalidHashLength)
			assert.Equal(suite.T(), found, false)

			found2, err2 := suite.service.CompareHashAndPassword(ctx, "wrong_password", nil)
			assert.Error(suite.T(), err2)
			assert.ErrorContains(suite.T(), err2, "invalid hash length")
			assert.ErrorIs(suite.T(), err2, ErrInvalidHashLength)
			assert.Equal(suite.T(), found2, false)
		})
		suite.Run("Password Mismatch", func() {
			hash, _ := suite.service.HashPassword("password")

			ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
			defer cancel()
			found, err := suite.service.CompareHashAndPassword(ctx, "wrong_password", hash)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "timeout exceeded")
			assert.ErrorIs(suite.T(), err, ErrPasswordMismatch)
			assert.Equal(suite.T(), found, false)
		})
	})
}
