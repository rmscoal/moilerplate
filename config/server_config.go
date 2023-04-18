package config

import (
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type serverConfig struct {
	Host string
	Port string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newServerConfig() {
	s := serverConfig{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
	}

	if err := s.validate(); err != nil {
		log.Fatalf("%s", err)
	}

	c.Server = s
}

// validate method    validates the serverConfig
// values such that it meets the requirements.
func (d serverConfig) validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Host, validation.Required, is.Host),
		validation.Field(&d.Port, validation.Required, is.Port),
	)
}
