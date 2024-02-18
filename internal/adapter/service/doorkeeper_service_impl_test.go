package service

import (
	"context"
	"testing"
	"time"

	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	"github.com/rmscoal/moilerplate/testing/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

type DoorkeeperServiceTestSuite struct {
	suite.Suite
	dk *doorkeeper.Doorkeeper
}

func TestDoorkeeperService(t *testing.T) {
	suite.Run(t, new(DoorkeeperServiceTestSuite))
}

func (suite *DoorkeeperServiceTestSuite) SetupSuite() {
	otel.SetTracerProvider(observability.NewTestingTracerProvider())

	dk := doorkeeper.GetDoorkeeper(
		// JWT
		doorkeeper.RegisterJWTIssuer("TESTING"),
		doorkeeper.RegisterJWTSignMethod("HMAC", "256"),
		doorkeeper.RegisterJWTPublicKey("verysecretkey"),
		doorkeeper.RegisterJWTPrivateKey("verysecretkey"),
		doorkeeper.RegisterJWTAccessDuration(5*time.Minute),
		doorkeeper.RegisterJWTRefreshDuration(10*time.Minute),
		// Encryptor
		doorkeeper.RegisterEncryptorSecretKey("verystrongsecretkey"),
		// General
		doorkeeper.RegisterGeneralHasherFunc("SHA384"),
	)

	suite.dk = dk
}

func (suite *DoorkeeperServiceTestSuite) SetupTest() {}

func (suite *DoorkeeperServiceTestSuite) TestHashAndEncodeStringWithSalt_Success() {
	service := NewDoorkeeperService(suite.dk)

	result := service.HashAndEncodeStringWithSalt(context.Background(), "password", "salt")
	assert.NotEmpty(suite.T(), result)
}
