package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type appConfig struct {
	Environment           string
	LogPath               string
	DefaultPaginationSize int

	defaultPaginationSizeStr string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newAppConfig() {
	a := appConfig{
		Environment:              strings.ToUpper(os.Getenv("ENVIRONMENT")),
		LogPath:                  strings.ToLower(os.Getenv("LOG_PATH")),
		defaultPaginationSizeStr: os.Getenv("DEFAULT_ROWS_PER_PAGE"),
	}

	if err := a.validate(); err != nil {
		log.Fatalf("FATAL - %s", err)
	}

	b := &a
	b.parse()

	c.App = *b
}

// validate method    validates the serverConfig
// values such that it meets the requirements.
func (a appConfig) validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.defaultPaginationSizeStr, validation.Required, is.Int),
		validation.Field(&a.LogPath, validation.Required),
		validation.Field(&a.Environment, validation.Required, validation.In(
			"MIGRATION",
			"DEVELOPMENT",
			"TESTING",
			"STAGING",
			"PRODUCTION",
		)),
	)
}

// parse method    parses a string value from env to
// the dedication type destination.
func (a *appConfig) parse() {
	defaultPaginationSize, err := strconv.Atoi(a.defaultPaginationSizeStr)
	if err != nil {
		log.Fatalf("%s", err)
	}

	a.DefaultPaginationSize = defaultPaginationSize
}
