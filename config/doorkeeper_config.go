package config

import (
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type doorkeeperConfig struct {
	HashSalt      string
	SecretKey     string
	SigningMethod string
	SignSize      string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newDoorkeeperConfig() {
	d := doorkeeperConfig{
		HashSalt:      os.Getenv("DOORKEEPER_HASH_SALT"),
		SecretKey:     os.Getenv("DOORKEEPER_SECRET_KEY"),
		SigningMethod: os.Getenv("DOORKEEPER_SIGNING_METHOD"),
		SignSize:      os.Getenv("DOORKEEPER_SIGN_SIZE"),
	}

	if err := d.validate(); err != nil {
		log.Fatalf("FATAL - %s", err)
	}

	c.Doorkeeper = d
}

// validate method    validates the serverConfig
// values such that it meets the requirements.
func (d doorkeeperConfig) validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.HashSalt, validation.Required, validation.Length(20, 10000)),
		validation.Field(&d.SecretKey, validation.Required, validation.Length(20, 10000)),
		validation.Field(&d.SigningMethod, validation.Required,
			validation.In("HMAC", "RSA", "ECDSA", "RSA-PSS", "EdDSA")),
		validation.Field(&d.SignSize, validation.Required,
			validation.When(d.SigningMethod != "EdDSA", validation.In("256, 384, 512"))),
	)
}
