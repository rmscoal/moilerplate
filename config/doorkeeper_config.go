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
	// --- JWT Config ---
	AccessDuration     time.Duration
	accessDurationStr  string
	RefreshDuration    time.Duration
	refreshDurationStr string

	SigningMethod string
	SignSize      string
	Path          string
	Issuer        string

	// --- Password Hasher Config ---
	HashMethod string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newDoorkeeperConfig() {
	d := doorkeeperConfig{
		accessDurationStr:  strings.ToLower(os.Getenv("DOORKEEPER_ACCESS_TOKEN_DURATION")),
		refreshDurationStr: strings.ToLower(os.Getenv("DOORKEEPER_REFRESH_TOKEN_DURATION")),
		SigningMethod:      strings.ToUpper(os.Getenv("DOORKEEPER_SIGNING_METHOD")),
		SignSize:           strings.ToLower(os.Getenv("DOORKEEPER_SIGN_SIZE")),
		Path:               strings.ToLower(os.Getenv("DOORKEEPER_CERT_PATH")),
		Issuer:             strings.ToUpper(os.Getenv("DOORKEEPER_ISSUER")),
		HashMethod:         strings.ToUpper(os.Getenv("DOORKEEPER_HASH_METHOD")),
	}

	if err := d.validate(); err != nil {
		log.Fatalf("FATAL - %s", err)
	}

	if d.accessDurationStr != "" {
		d.AccessDuration = d.parseTime(d.accessDurationStr)
	}

	if d.refreshDurationStr != "" {
		d.RefreshDuration = d.parseTime(d.refreshDurationStr)
	}

	c.Doorkeeper = d
}

// validate method    validates the serverConfig
// values such that it meets the requirements.
func (d doorkeeperConfig) validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.SigningMethod, validation.Required.
			Error("Please provide a signing method in the environment. This is needed for signing authorization tokens"),
			validation.In("HMAC", "RSA", "ECDSA", "RSA-PSS", "EdDSA")),
		validation.Field(&d.SignSize,
			validation.When(d.SigningMethod != "EdDSA", validation.Required.
				Error("Please provide a signing size in the environment. This is needed for signing authorization tokens"),
				validation.In("256", "384", "512"))),
		validation.Field(&d.HashMethod, validation.Required.
			Error("Please provide hash method in the environment. This is needed when hashing credentials"),
			validation.In("SHA1", "SHA224", "SHA256",
				"SHA384", "SHA512", "SHA3_224", "SHA3_256", "SHA3_384", "SHA3_512")),
		validation.Field(&d.accessDurationStr, validation.When(d.accessDurationStr != "", validation.By(d.validateTime))),
		validation.Field(&d.Path, validation.When(d.SigningMethod != "HMAC", validation.Required)),
		validation.Field(&d.Issuer, validation.Required),
	)
}

func (d doorkeeperConfig) validateTime(value any) error {
	timeSlc := strings.Split(value.(string), " ")
	// Checks length
	if len(timeSlc) != 2 {
		return fmt.Errorf("required length is 2 but got %d", len(timeSlc))
	}
	// Checks time meter
	if err := validation.Validate(&timeSlc[1],
		validation.In("second", "seconds", "minute", "minutes", "hour", "hours").
			Error("invalid time unit")); err != nil {
		return err
	}
	// Check time value
	_, err := strconv.Atoi(timeSlc[0])
	if err != nil {
		return err
	}
	return nil
}

func (d doorkeeperConfig) parseTime(str string) time.Duration {
	var res time.Duration
	timeSlc := strings.Split(str, " ")
	switch timeSlc[1] {
	case "second", "seconds":
		res = time.Second
	case "minute", "minutes":
		res = time.Minute
	case "hour", "hours":
		res = time.Hour
	}
	// Check time value
	val, _ := strconv.Atoi(timeSlc[0])
	res = time.Duration(val) * res
	return res
}
