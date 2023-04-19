package config

import (
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type doorkeeperConfig struct {
	hashSalt      string
	hashMethod    string
	secretKey     string
	signingMethod string
	signSize      string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newDoorkeeperConfig() {
	d := doorkeeperConfig{
		hashSalt:      os.Getenv("DOORKEEPER_HASH_SALT"),
		hashMethod:    os.Getenv("DOORKEEPER_HASH_METHOD"),
		secretKey:     os.Getenv("DOORKEEPER_SECRET_KEY"),
		signingMethod: os.Getenv("DOORKEEPER_SIGNING_METHOD"),
		signSize:      os.Getenv("DOORKEEPER_SIGN_SIZE"),
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
		validation.Field(&d.hashSalt, validation.Required, validation.Length(20, 10000)),
		validation.Field(&d.secretKey, validation.Required, validation.Length(20, 10000)),
		validation.Field(&d.signingMethod, validation.Required,
			validation.In("HMAC", "RSA", "ECDSA", "RSA-PSS", "EdDSA")),
		validation.Field(&d.signSize, validation.Required,
			validation.When(d.signingMethod != "EdDSA", validation.In("256, 384, 512"))),
		validation.Field(&d.hashMethod, validation.Required,
			validation.In("MD4", "MD5", "SHA1", "SHA224", "SHA256",
				"SHA384", "SHA512", "SHA3_224", "SHA3_256", "SHA3_384", "SHA3_512")),
	)
}

func (d *doorkeeperConfig) HashSalt() string {
	return d.hashSalt
}

func (d *doorkeeperConfig) HashMethod() string {
	return d.hashMethod
}

func (d *doorkeeperConfig) SecretKey() string {
	return d.secretKey
}

func (d *doorkeeperConfig) SigningMethod() string {
	return d.signingMethod
}

func (d *doorkeeperConfig) SigningSize() string {
	return d.signSize
}
