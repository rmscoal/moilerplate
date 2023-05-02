package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/service"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/doorkeeper"
	"github.com/stretchr/testify/suite"
)

var passwordTestCases = []string{"password", "mediumpassword", "verylongpassword", "nvsjvn&#*@jdsnvvavakac#@HBSDVs"}

type DoorkeeperServiceImplBenchSuite struct {
	suite.Suite
	dk      *doorkeeper.Doorkeeper
	service service.IDoorkeeperService
}

func (suite *DoorkeeperServiceImplBenchSuite) SetupTest() {
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

func (suite *DoorkeeperServiceImplBenchSuite) SetupSuite() {}

func (suite *DoorkeeperServiceImplBenchSuite) TearDownTest() {
	suite.dk = nil
	suite.service = nil
}

func (suite *DoorkeeperServiceImplBenchSuite) TearDownSuite() {}

// go test ./internal/adapter/service -bench=Benchmark -count=5 -benchtime=100x -benchmem -run=^#
func BenchmarkHashPassword(b *testing.B) {
	s := new(DoorkeeperServiceImplBenchSuite)
	s.SetT(&testing.T{})
	s.SetupSuite()
	for _, password := range passwordTestCases {
		b.ResetTimer()
		b.Run(fmt.Sprintf("password_of_%s", password), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.SetupTest()
				b.StartTimer()

				s.service.HashPassword(password)

				b.StopTimer()
				s.TearDownTest()
			}
		})
	}
	s.TearDownSuite()
}
