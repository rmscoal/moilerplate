package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type doorkeeperConfig struct {
	// --- JWT Config ---
	jwtIssuer          string
	jwtSignMethod      string
	jwtSignSize        string
	jwtPrivateKey      string
	jwtPublicKey       string
	jwtAccessDuration  time.Duration
	jwtRefreshDuration time.Duration

	// --- General ---
	generalHashMethod string

	// --- Encryptor ---
	encryptorSecretKey string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newDoorkeeperConfig() {
	d := doorkeeperConfig{
		// JWT
		jwtIssuer:     strings.ToUpper(os.Getenv("DOORKEEPER_JWT_ISSUER")),
		jwtSignMethod: strings.ToUpper(os.Getenv("DOORKEEPER_JWT_SIGNING_METHOD")),
		jwtSignSize:   os.Getenv("DOORKEEPER_JWT_SIGN_SIZE"),
		jwtPrivateKey: os.Getenv("DOORKEEPER_JWT_PRIV_KEY"),
		jwtPublicKey:  os.Getenv("DOORKEEPER_JWT_PUB_KEY"),
		// General
		generalHashMethod: strings.ToUpper(os.Getenv("DOORKEEPER_GENERAL_HASH_METHOD")),
		// Encryptor
		encryptorSecretKey: os.Getenv("DOORKEEPER_ENCRYPTOR_SECRET_KEY"),
	}
	if err := d.parse(); err != nil {
		log.Fatalf("Error parsing doorkeeper environment: %s\n", err)
	}

	if err := d.validate(); err != nil {
		log.Fatalf("FATAL - %s", err)
	}

	c.Doorkeeper = d
}

func (d doorkeeperConfig) validate() error {
	return validation.ValidateStruct(&d,
		// JWT
		validation.Field(&d.jwtIssuer, validation.Required),
		validation.Field(&d.jwtSignMethod, validation.Required.
			Error("Please provide a signing method in the environment. This is needed for signing authorization tokens"),
			validation.In("HMAC", "RSA", "ECDSA", "RSA-PSS", "EdDSA")),
		validation.Field(&d.jwtSignSize,
			validation.When(d.jwtSignMethod != "EdDSA", validation.Required.
				Error("Please provide a signing size in the environment. This is needed for signing authorization tokens"),
				validation.In("256", "384", "512"))),
		validation.Field(&d.jwtPrivateKey, validation.Required, validation.Length(10, 0),
			validation.When(d.jwtSignMethod == "HMAC", validation.In(d.jwtPublicKey))),
		validation.Field(&d.jwtPublicKey, validation.Required, validation.Length(10, 0),
			validation.When(d.jwtSignMethod == "HMAC", validation.In(d.jwtPrivateKey)),
		),

		// General
		validation.Field(&d.generalHashMethod, validation.Required.
			Error("Please provide hash method in the environment. This is needed when hashing credentials"),
			validation.In("SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "SHA3_224", "SHA3_256", "SHA3_384", "SHA3_512")),

		// Encryptor
		validation.Field(&d.encryptorSecretKey, validation.Required, validation.Length(10, 0)),
	)
}

// parse method    parses a string value from env to
// the dedication type destination.
func (d *doorkeeperConfig) parse() (err error) {
	d.jwtAccessDuration, err = time.ParseDuration(os.Getenv("DOORKEEPER_JWT_ACCESS_TOKEN_DURATION"))
	if err != nil {
		return fmt.Errorf("unable to parse DOORKEEPER_JWT_ACCESS_TOKEN_DURATION: %s", err)
	}

	d.jwtRefreshDuration, err = time.ParseDuration(os.Getenv("DOORKEEPER_JWT_REFRESH_TOKEN_DURATION"))
	if err != nil {
		return fmt.Errorf("unable to parse DOORKEEPER_JWT_REFRESH_TOKEN_DURATION: %s", err)
	}

	return nil
}

// Getter functions for getting doorkeeper configurations
// --- JWT ---

func (d doorkeeperConfig) JWTSigningMethod() string {
	return d.jwtSignMethod
}

func (d doorkeeperConfig) JWTSignSize() string {
	return d.jwtSignSize
}

func (d doorkeeperConfig) JWTIssuer() string {
	return d.jwtIssuer
}

func (d doorkeeperConfig) JWTPublicKey() string {
	return d.jwtPublicKey
}

func (d doorkeeperConfig) JWTPrivateKey() string {
	return d.jwtPrivateKey
}

func (d doorkeeperConfig) JWTAccessTokenDuration() time.Duration {
	return d.jwtAccessDuration
}

func (d doorkeeperConfig) JWTRefreshTokenDuration() time.Duration {
	return d.jwtRefreshDuration
}

// --- General ---

func (d doorkeeperConfig) GeneralHashMethod() string {
	return d.generalHashMethod
}

// --- Encryptor ---

func (d doorkeeperConfig) EncryptorSecretKey() string {
	return d.encryptorSecretKey
}
