package service

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/service"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/doorkeeper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	// Context for testing
	TEST_CONTEXT = context.Background()

	// Dummy user domain for testing
	USER_DOMAIN = domain.User{
		Id: "DOMAIN_USER_ID",
		Credential: vo.UserCredential{Tokens: vo.UserToken{
			TokenID:  "TOKEN_ID",
			IssuedAt: time.Now(),
		}},
	}
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
		// This test seems to fluctuate ⬆️
	})

	suite.Run("Panic on Generating Random Range", func() {
		// It seems that HashPassword will only throw error iff
		// MinSaltLength > MaxSaltLength or rand.Reader fails on
		// the os side.
		MinSaltLength, MaxSaltLength = MaxSaltLength, MinSaltLength

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

func (suite *DoorkeeperServiceImplTestSuite) TestGenerateTokens() {
	suite.Run("Successful Generate Tokens", func() {
		userToken, err := suite.service.GenerateUserTokens(USER_DOMAIN)
		assert.Nil(suite.T(), err)
		assert.NotEmpty(suite.T(), userToken)
		assert.NotEmpty(suite.T(), userToken.AccesssToken)
		assert.NotEmpty(suite.T(), userToken.RefreshToken)
		assert.Equal(suite.T(), false, userToken.Issued)
	})
}

func (suite *DoorkeeperServiceImplTestSuite) TestVerifyParseAccessToken() {
	suite.Run("Successful Verify and Parse", func() {
		// Generate userToken
		userToken, _ := suite.service.GenerateUserTokens(USER_DOMAIN)

		// Start verify and parse test
		userId, err := suite.service.VerifyAndParseToken(TEST_CONTEXT, userToken.AccesssToken)
		assert.Equal(suite.T(), USER_DOMAIN.Id, userId)
		assert.Nil(suite.T(), err)
	})

	suite.Run("Unsuccessful Verify and Parse", func() {
		suite.Run("Invalid Access Token", func() {
			accessToken := "INVALID_ACCESS_TOKEN"

			// Start verify and parse test
			userId, err := suite.service.VerifyAndParseToken(TEST_CONTEXT, accessToken)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "validation failed")
			assert.Empty(suite.T(), userId)
		})

		suite.Run("Invalid Signing Method", func() {
			// This is a PS256 Signed Dummy JWT from jwt.io
			userToken := "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.iOeNU4dAFFeBwNj6qdhdvm-IvDQrTa6R22lQVJVuWJxorJfeQww5Nwsra0PjaOYhAMj9jNMO5YLmud8U7iQ5gJK2zYyepeSuXhfSi8yjFZfRiSkelqSkU19I-Ja8aQBDbqXf2SAWA8mHF8VS3F08rgEaLCyv98fLLH4vSvsJGf6ueZSLKDVXz24rZRXGWtYYk_OYYTVgR1cg0BLCsuCvqZvHleImJKiWmtS0-CymMO4MMjCy_FIl6I56NqLE9C87tUVpo1mT-kbg5cHDD8I7MjCW5Iii5dethB4Vid3mZ6emKjVYgXrtkOQ-JyGMh6fnQxEFN1ft33GX2eRHluK9eg"

			userId, err := suite.service.VerifyAndParseToken(TEST_CONTEXT, userToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "signing method invalid")

			// This is a HS256 Signed Dummy JWT from jwt.io
			userToken = "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.iOeNU4dAFFeBwNj6qdhdvm-IvDQrTa6R22lQVJVuWJxorJfeQww5Nwsra0PjaOYhAMj9jNMO5YLmud8U7iQ5gJK2zYyepeSuXhfSi8yjFZfRiSkelqSkU19I-Ja8aQBDbqXf2SAWA8mHF8VS3F08rgEaLCyv98fLLH4vSvsJGf6ueZSLKDVXz24rZRXGWtYYk_OYYTVgR1cg0BLCsuCvqZvHleImJKiWmtS0-CymMO4MMjCy_FIl6I56NqLE9C87tUVpo1mT-kbg5cHDD8I7MjCW5Iii5dethB4Vid3mZ6emKjVYgXrtkOQ-JyGMh6fnQxEFN1ft33GX2eRHluK9eg"

			userId, err = suite.service.VerifyAndParseToken(TEST_CONTEXT, userToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "signing method invalid")
		})

		suite.Run("Expired Access Token", func() {
			// This is a RS256 Signed Valid Expired JWT
			userToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE2ODI5MjgzMzQsImlhdCI6MTY4MjkyODAzNCwiaXNzIjoiTU9JTE9SUExBVEUiLCJuYmYiOjE2ODI5MjgwMzQsInVzZXJJZCI6ImNmNzZhZWY3LWIzYTUtNDVlYy04MmJiLTMwYjhhYWUzN2M0NSJ9.pollns246zsejpZZUlaDghu-j6A0bnzTh8G_lQr6fFh39mykM3QrCoYrhy07BRPbsSHrY2w_qrNlEPjFWoVS3WzfM6z-c60RzOKOmEZ7P98YgC67eCSONAwRabEj5QFbrIRACsrLF0OBCjM0alDIUMmon8IG3xD5uioZPyBR4S4bM5fcb5f85xbUF7vQOgbWSr9b_5LcT980CwRAK9957NBCAoGWBG0Y0X1aCoPUK_u5U5nPRx8oykHDCcGqPs44xECRSxQ6oSvrTn_oItWKJYty7_tMAdgHEpIiBgIUrcQiTIpDOhwUbLjBrMyR4xGNv8UusyLJgk0d5H5Y_ERbvw"

			userId, err := suite.service.VerifyAndParseToken(TEST_CONTEXT, userToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "token has expired")
		})

		suite.Run("Invalid Claims", func() {
			// This is a RS256 Signed Valid Expired Refresh JWT
			userToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE2ODI5MzE2MzQsImlhdCI6MTY4MjkyODAzNCwiaXNzIjoiTU9JTE9SUExBVEUiLCJqdGkiOiJhZThlNDkxYi1mNDZiLTQ4YzItYjc0ZC00ODYzYTc4ODg5ODUiLCJuYmYiOjE2ODI5MjgwMzR9.nNJXT7Bq1mxqLb0uPTpjqA3cK3YTu7FZuQNJiEaEE1Wkm-OiEvUYYboXilxuc49Pm3yohC-C82NdBw5hUZt0owb7SCK5R5NMij3wlD4rvUKOQ5mFHkdROB4IAuJiEP3qJ8xkSFrCSYCxLOhK-P85d5VBLvM07ptakL3AUXLan7z_GtaQ7P8wlb0zJs5ElsoBWRQB36znkZhn8qRTSWUYHsoh9QV7OiVLS7gS6rFrfKHa9bzeiwbiqmUi7ivnptgW4mZlDG7aCSzI56RYh3dt3tpnvyciyN9me7eQbr3WjwzlRAmNMnJyfoKr5IcsZS3MzRWY_BSJz4Zc7Cb4LmTuWg"

			userId, err := suite.service.VerifyAndParseToken(TEST_CONTEXT, userToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "does not contained required claim")
		})
	})
}

func (suite *DoorkeeperServiceImplTestSuite) TestVerifyParseRefreshToken() {
	suite.Run("Successful Verify and Parse", func() {
		// Generate userToken
		userToken, _ := suite.service.GenerateUserTokens(USER_DOMAIN)

		// Start verify and parse test
		jti, err := suite.service.VerifyAndParseRefreshToken(TEST_CONTEXT, userToken.RefreshToken)
		assert.Equal(suite.T(), USER_DOMAIN.Credential.Tokens.TokenID, jti)
		assert.Nil(suite.T(), err)
	})

	suite.Run("Unsuccessful Verify and Parse", func() {
		// Generate userToken
		generatedUserToken, _ := suite.service.GenerateUserTokens(USER_DOMAIN)

		suite.Run("Invalid Access Token", func() {
			accessToken := "INVALID_ACCESS_TOKEN"

			// Start verify and parse test
			userId, err := suite.service.VerifyAndParseRefreshToken(TEST_CONTEXT, accessToken)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "validation failed")
			assert.Empty(suite.T(), userId)
		})
		suite.Run("Invalid Signing Method", func() {
			// This is a PS256 Signed Dummy JWT from jwt.io
			userToken := "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.iOeNU4dAFFeBwNj6qdhdvm-IvDQrTa6R22lQVJVuWJxorJfeQww5Nwsra0PjaOYhAMj9jNMO5YLmud8U7iQ5gJK2zYyepeSuXhfSi8yjFZfRiSkelqSkU19I-Ja8aQBDbqXf2SAWA8mHF8VS3F08rgEaLCyv98fLLH4vSvsJGf6ueZSLKDVXz24rZRXGWtYYk_OYYTVgR1cg0BLCsuCvqZvHleImJKiWmtS0-CymMO4MMjCy_FIl6I56NqLE9C87tUVpo1mT-kbg5cHDD8I7MjCW5Iii5dethB4Vid3mZ6emKjVYgXrtkOQ-JyGMh6fnQxEFN1ft33GX2eRHluK9eg"

			userId, err := suite.service.VerifyAndParseRefreshToken(TEST_CONTEXT, userToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "signing method invalid")

			// This is a HS256 Signed Dummy JWT from jwt.io
			userToken = "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.iOeNU4dAFFeBwNj6qdhdvm-IvDQrTa6R22lQVJVuWJxorJfeQww5Nwsra0PjaOYhAMj9jNMO5YLmud8U7iQ5gJK2zYyepeSuXhfSi8yjFZfRiSkelqSkU19I-Ja8aQBDbqXf2SAWA8mHF8VS3F08rgEaLCyv98fLLH4vSvsJGf6ueZSLKDVXz24rZRXGWtYYk_OYYTVgR1cg0BLCsuCvqZvHleImJKiWmtS0-CymMO4MMjCy_FIl6I56NqLE9C87tUVpo1mT-kbg5cHDD8I7MjCW5Iii5dethB4Vid3mZ6emKjVYgXrtkOQ-JyGMh6fnQxEFN1ft33GX2eRHluK9eg"

			userId, err = suite.service.VerifyAndParseRefreshToken(TEST_CONTEXT, userToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "signing method invalid")
		})

		suite.Run("Expired Refresh Token", func() {
			// This is a RS256 Signed Valid Expired Refresh JWT
			userToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlYXQiOjE2ODI5MzE2MzQsImlhdCI6MTY4MjkyODAzNCwiaXNzIjoiTU9JTE9SUExBVEUiLCJqdGkiOiJhZThlNDkxYi1mNDZiLTQ4YzItYjc0ZC00ODYzYTc4ODg5ODUiLCJuYmYiOjE2ODI5MjgwMzR9.nNJXT7Bq1mxqLb0uPTpjqA3cK3YTu7FZuQNJiEaEE1Wkm-OiEvUYYboXilxuc49Pm3yohC-C82NdBw5hUZt0owb7SCK5R5NMij3wlD4rvUKOQ5mFHkdROB4IAuJiEP3qJ8xkSFrCSYCxLOhK-P85d5VBLvM07ptakL3AUXLan7z_GtaQ7P8wlb0zJs5ElsoBWRQB36znkZhn8qRTSWUYHsoh9QV7OiVLS7gS6rFrfKHa9bzeiwbiqmUi7ivnptgW4mZlDG7aCSzI56RYh3dt3tpnvyciyN9me7eQbr3WjwzlRAmNMnJyfoKr5IcsZS3MzRWY_BSJz4Zc7Cb4LmTuWg"

			userId, err := suite.service.VerifyAndParseRefreshToken(TEST_CONTEXT, userToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "token has expired")
		})

		suite.Run("Invalid Claims", func() {
			userId, err := suite.service.VerifyAndParseRefreshToken(TEST_CONTEXT, generatedUserToken.AccesssToken)
			assert.Empty(suite.T(), userId)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "does not contained required claim")
		})
	})
}
