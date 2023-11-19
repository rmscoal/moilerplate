package config

import (
	"log"
	"os"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type doorkeeperConfig struct {
	// --- JWT Config ---
	accessDuration  time.Duration
	refreshDuration time.Duration

	signingMethod string
	signSize      string
	secretKey     string
	certPath      string
	issuer        string

	// --- Password Hasher Config ---
	hashMethod string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newDoorkeeperConfig() {
	d := doorkeeperConfig{
		signingMethod: strings.ToUpper(os.Getenv("DOORKEEPER_SIGNING_METHOD")),
		signSize:      strings.ToLower(os.Getenv("DOORKEEPER_SIGN_SIZE")),
		certPath:      strings.ToLower(os.Getenv("DOORKEEPER_CERT_PATH")),
		issuer:        strings.ToUpper(os.Getenv("DOORKEEPER_ISSUER")),
		hashMethod:    strings.ToUpper(os.Getenv("DOORKEEPER_HASH_METHOD")),
		secretKey:     os.Getenv("DOORKEEPER_SECRET_KEY"),
	}
	if err := (&d).parse(); err != nil {
		log.Fatalf("Error parsing doorkeeper environment: %s\n", err)
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
		validation.Field(&d.signingMethod, validation.Required.
			Error("Please provide a signing method in the environment. This is needed for signing authorization tokens"),
			validation.In("HMAC", "RSA", "ECDSA", "RSA-PSS", "EdDSA")),
		validation.Field(&d.signSize,
			validation.When(d.signingMethod != "EdDSA", validation.Required.
				Error("Please provide a signing size in the environment. This is needed for signing authorization tokens"),
				validation.In("256", "384", "512"))),
		validation.Field(&d.hashMethod, validation.Required.
			Error("Please provide hash method in the environment. This is needed when hashing credentials"),
			validation.In("SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "SHA3_224", "SHA3_256", "SHA3_384", "SHA3_512")),
		validation.Field(&d.certPath, validation.When(d.signingMethod != "HMAC", validation.Required)),
		validation.Field(&d.secretKey, validation.When(d.signingMethod == "HMAC", validation.Required.
			Error("a HMAC signing method requires a secret key"))),
		validation.Field(&d.issuer, validation.Required),
	)
}

// parse method    parses a string value from env to
// the dedication type destination.
func (d *doorkeeperConfig) parse() (err error) {
	d.accessDuration, err = time.ParseDuration(os.Getenv("DOORKEEPER_ACCESS_TOKEN_DURATION"))
	if err != nil {
		return err
	}

	d.refreshDuration, err = time.ParseDuration(os.Getenv("DOORKEEPER_REFRESH_TOKEN_DURATION"))
	if err != nil {
		return err
	}

	return nil
}

// Getter functions for getting doorkeeper configurations
func (d doorkeeperConfig) HashMethod() string {
	return d.hashMethod
}

func (d doorkeeperConfig) SigningMethod() string {
	return d.signingMethod
}

func (d doorkeeperConfig) SignSize() string {
	return d.signSize
}

func (d doorkeeperConfig) Issuer() string {
	return d.issuer
}

func (d doorkeeperConfig) CertPath() string {
	return d.certPath
}

func (d doorkeeperConfig) SecretKey() string {
	return d.secretKey
}

func (d doorkeeperConfig) AccessTokenDuration() time.Duration {
	return d.accessDuration
}

func (d doorkeeperConfig) RefreshTokenDuration() time.Duration {
	return d.refreshDuration
}
