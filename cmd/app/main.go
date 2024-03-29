package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
	"github.com/mitchellh/cli"
	"github.com/rmscoal/moilerplate/config"
	"github.com/rmscoal/moilerplate/internal/composer"
	v1 "github.com/rmscoal/moilerplate/internal/delivery/v1"
	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	httpserver "github.com/rmscoal/moilerplate/pkg/http"
	"github.com/rmscoal/moilerplate/pkg/observability"
	"github.com/rmscoal/moilerplate/pkg/postgres"
	"github.com/rmscoal/moilerplate/pkg/rater"
	"github.com/rmscoal/moilerplate/swagger"
)

type app struct {
	flagSecure                            bool
	flagEnvPath                           string
	flagServerCertPath, flagServerKeyPath string
	flagMode                              string
}

func NewAppCLI() *app {
	return &app{}
}

func (a *app) Synopsis() string {
	return "runs the server"
}

func (a *app) Help() string {
	return `
	Usage: server [--with-secure] [--help | -h] [--cert=<cert_path>] [--key=<key_path>] [--env-path=<env_path>]
	              [--mode=<DEVELOPMENT|TESTING|PRODUCTION>]
	Spins up the server...

		--mode            start the server in either DEVELOPMENT, TESTING, or PRODUCTION mode
	
	If you want to start the server with TLS enabled, these flags might be useful:
		--with-secure     with-secure will start the server using TLS enabled
		--cert            provide with the certificate path
		--key             provide with the key path
	
	If you want to provide the app to read from provided dot env files, these flags might be useful:
		--env-path        provide with the dot env file path
	`
}

func main() {
	cli := cli.NewCLI("moilerplate", "1.0.0")
	cli.Commands = commands()
	cli.Args = os.Args[1:]

	exitCode, err := cli.Run()
	if err != nil {
		log.Fatalf("unable to run app")
	}

	os.Exit(exitCode)
}

func commands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return NewAppCLI(), nil
		},
	}
}

func (a *app) flags() *flag.FlagSet {
	f := flag.NewFlagSet("server", flag.ExitOnError)

	f.BoolVar(&a.flagSecure, "with-secure", false, "with-secure will start the server in https with CA cert required")
	f.StringVar(&a.flagMode, "mode", "DEVELOPMENT", "mode indicates the mode of the app, there are DEVELOPMENT, TESTING and PRODUCTION")
	f.StringVar(&a.flagEnvPath, "env-path", "", "path to the environment variable to read from, for example '.env'")
	f.StringVar(&a.flagServerCertPath, "cert", "", "server CA certificate path for https")
	f.StringVar(&a.flagServerKeyPath, "key", "", "server key path for https")

	return f
}

func (a *app) Run(args []string) int {
	f := a.flags()
	if err := f.Parse(args); err != nil {
		log.Println("Parsing flag error", err)
		return 1
	}

	if err := a.validateFlags(); err != nil {
		log.Fatal(err)
	}

	if err := a.loadEnv(); err != nil {
		log.Fatal(err)
	}

	// Swagger documentation info
	swagger.SwaggerInfo.Title = "Moilerplate"
	swagger.SwaggerInfo.Description = "A monolithic RESTful API for Go"
	swagger.SwaggerInfo.Version = "1.0"
	swagger.SwaggerInfo.Host = "localhost:8082"
	swagger.SwaggerInfo.BasePath = "/api/v1"
	swagger.SwaggerInfo.Schemes = []string{"http", "https"}

	cfg := config.GetConfig()

	// Postgres .-.
	pg := postgres.GetPostgres(
		cfg.Db.URL,
		postgres.MaxPoolSize(cfg.Db.MaxPoolSize()),
		postgres.MaxOpenCoon(cfg.Db.MaxOpenConn()),
		postgres.SetMode(a.flagMode),
	)

	// Doorkeeper .-.
	dk := doorkeeper.GetDoorkeeper(
		// JWT
		doorkeeper.RegisterJWTIssuer(cfg.Doorkeeper.JWTIssuer()),
		doorkeeper.RegisterJWTSignMethod(cfg.Doorkeeper.JWTSigningMethod(), cfg.Doorkeeper.JWTSignSize()),
		doorkeeper.RegisterJWTPublicKey(cfg.Doorkeeper.JWTPublicKey()),
		doorkeeper.RegisterJWTPrivateKey(cfg.Doorkeeper.JWTPrivateKey()),
		doorkeeper.RegisterJWTRefreshDuration(cfg.Doorkeeper.JWTRefreshTokenDuration()),
		doorkeeper.RegisterJWTAccessDuration(cfg.Doorkeeper.JWTAccessTokenDuration()),
		// General
		doorkeeper.RegisterGeneralHasherFunc(cfg.Doorkeeper.GeneralHashMethod()),
		// Encryptor
		doorkeeper.RegisterEncryptorSecretKey(cfg.Doorkeeper.EncryptorSecretKey()),
	)

	// Rate Limitter .-.
	rt := rater.GetRater(context.Background(),
		rater.RegisterRateLimitForEachClient(cfg.App.RaterLimit()),
		rater.RegisterBurstLimitForEachClient(cfg.App.BurstLimit()),
		rater.RegisterEvaluationInterval(cfg.App.RaterEvaluationInterval()),
		rater.RegisterDeletionTime(cfg.App.RaterDeletionTime()),
	)

	// Opentelemetry
	observability.Init(context.Background(),
		observability.TraceEndpoint(cfg.Otel.GetTraceEndpoint()),
		observability.MetricsEndpoint(cfg.Otel.GetMetricEndpoint()),
		observability.ServiceName(cfg.Otel.GetServiceName()),
		observability.ServiceVersion(cfg.Otel.GetServiceVersion()),
		observability.ServiceInstanceID(cfg.Otel.GetServiceInstanceID()),
	)

	// Composers .-.
	serviceComposer := composer.NewServiceComposer(dk, rt, pg.Pool)
	repoComposer := composer.NewRepoComposer(pg)
	usecaseComposer := composer.NewUseCaseComposer(repoComposer, serviceComposer)

	deliveree := a.newDeliveryEngine()

	// Http
	v1.NewRouter(deliveree, usecaseComposer, serviceComposer)
	httpserver.NewServer(deliveree,
		httpserver.RegisterHostAndPort(cfg.Server.Host, cfg.Server.Port),
		httpserver.StartSecure(a.flagSecure, a.flagServerCertPath, a.flagServerKeyPath),
	)

	return 0
}

// loadEnv loads the environment
func (a *app) loadEnv() error {
	if a.flagEnvPath != "" {
		err := godotenv.Load(a.flagEnvPath)
		if err != nil {
			return fmt.Errorf("unable to load environment variable: %s", err)
		}
	}
	return nil
}

// newDeliveryEngine creates the new gin.engine depending on
// the environment state.
func (a *app) newDeliveryEngine() *gin.Engine {
	var deliveree *gin.Engine
	switch a.flagMode {
	case "PRODUCTION", "TESTING":
		gin.SetMode(gin.ReleaseMode)
		deliveree = gin.New()
		deliveree.Use(gin.Recovery())
	default:
		deliveree = gin.Default()
	}

	return deliveree
}

func (a app) validateFlags() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.flagServerCertPath, validation.When(a.flagSecure, validation.Required.Error("required to provide the server certificate path"))),
		validation.Field(&a.flagServerKeyPath, validation.When(a.flagSecure, validation.Required.Error("required to provide the server key path"))),
		validation.Field(&a.flagMode, validation.Required, validation.In("DEVELOPMENT", "TESTING", "PRODUCTION")),
	)
}
