package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type doorkeeperConfig struct {
	hashSalt      string
	hashMethod    string
	secretKey     string
	signingMethod string
	signSize      string

	Duration    time.Duration
	DurationStr string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newDoorkeeperConfig() {
	d := doorkeeperConfig{
		DurationStr:   strings.ToLower(os.Getenv("DOORKEEPER_TOKEN_DURATION")),
		hashSalt:      os.Getenv("DOORKEEPER_HASH_SALT"),
		hashMethod:    strings.ToUpper(os.Getenv("DOORKEEPER_HASH_METHOD")),
		secretKey:     os.Getenv("DOORKEEPER_SECRET_KEY"),
		signingMethod: strings.ToUpper(os.Getenv("DOORKEEPER_SIGNING_METHOD")),
		signSize:      os.Getenv("DOORKEEPER_SIGN_SIZE"),
	}

	if err := d.validate(); err != nil {
		log.Fatalf("FATAL - %s", err)
	}

	if d.DurationStr != "" {
		d.Duration = d.parseTime()
	}
	c.Doorkeeper = d
}

// validate method    validates the serverConfig
// values such that it meets the requirements.
func (d doorkeeperConfig) validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.hashSalt, validation.Required.
			Error("Please provide a hash salt in the environment. This is needed when hashing credentials"),
			validation.Length(20, 10000)),
		validation.Field(&d.secretKey,
			validation.When(d.signingMethod == "RSA", validation.Skip).Else(
				validation.Required.
					Error("Please provide a secret key in the environment. This is needed for signing authorization tokens"),
				validation.Length(20, 10000))),
		validation.Field(&d.signingMethod, validation.Required.
			Error("Please provide a signing method in the environment. This is needed for signing authorization tokens"),
			validation.In("HMAC", "RSA", "ECDSA", "RSA-PSS", "EdDSA")),
		validation.Field(&d.signSize,
			validation.When(d.signingMethod != "EdDSA", validation.Required.
				Error("Please provide a signing size in the environment. This is needed for signing authorization tokens"),
				validation.In("256", "384", "512"))),
		validation.Field(&d.hashMethod, validation.Required.
			Error("Please provide hash method in the environment. This is needed when hashing credentials"),
			validation.In("MD4", "MD5", "SHA1", "SHA224", "SHA256",
				"SHA384", "SHA512", "SHA3_224", "SHA3_256", "SHA3_384", "SHA3_512")),
		validation.Field(&d.Duration, validation.When(d.DurationStr != "",
			validation.By(
				func(value interface{}) error {
					timeSlc := strings.Split(value.(string), " ")
					// Checks length
					if len(timeSlc) != 2 {
						return fmt.Errorf("required length is 2 but got %d", len(timeSlc))
					}
					// Checks time meter
					if err := validation.Validate(&timeSlc[0],
						validation.In("second", "seconds", "minute", "minutes", "hour", "hours").
							Error("invalid time meter option")); err != nil {
						return err
					}
					// Check time value
					_, err := strconv.Atoi(timeSlc[1])
					if err != nil {
						return err
					}
					return nil
				},
			))),
	)
}

func (d doorkeeperConfig) parseTime() time.Duration {
	var res time.Duration
	timeSlc := strings.Split(d.DurationStr, " ")
	switch timeSlc[1] {
	case "second", "seconds":
		res = time.Second
	case "minute", "minutes":
		res = time.Minute
	case "hour", "hours":
		res = time.Hour
	}
	// Check time value
	val, _ := strconv.Atoi(timeSlc[1])
	res = time.Duration(val) * res
	return res
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
