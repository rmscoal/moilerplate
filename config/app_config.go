package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type appConfig struct {
	Environment             string
	LogPath                 string
	RaterLimit              int
	BurstLimit              int
	RaterEvaluationInterval time.Duration
	RaterDeletionTime       time.Duration
	DefaultPaginationSize   int

	defaultPaginationSizeStr string
	raterEvIntStr            string
	raterDelTime             string
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newAppConfig() {
	a := appConfig{
		Environment:              strings.ToUpper(os.Getenv("ENVIRONMENT")),
		LogPath:                  strings.ToLower(os.Getenv("LOG_PATH")),
		defaultPaginationSizeStr: os.Getenv("DEFAULT_ROWS_PER_PAGE"),
		raterEvIntStr:            os.Getenv("RATER_EVALUATION_INTERVAL"),
		raterDelTime:             os.Getenv("RATER_DELETION_TIME"),
	}

	b := &a
	b.parse()

	if err := a.validate(); err != nil {
		log.Fatalf("FATAL - %s", err)
	}

	a.RaterEvaluationInterval = a.parseTime(os.Getenv("RATER_EVALUATION_INTERVAL"))
	a.RaterDeletionTime = a.parseTime(os.Getenv("RATER_DELETION_TIME"))

	c.App = *b
}

// validate method    validates the serverConfig
// values such that it meets the requirements.
func (a appConfig) validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.defaultPaginationSizeStr, validation.Required, is.Int),
		validation.Field(&a.LogPath, validation.Required),
		validation.Field(&a.raterEvIntStr, validation.When(a.raterEvIntStr != "", validation.By(a.validateTime))),
		validation.Field(&a.raterDelTime, validation.When(a.raterDelTime != "", validation.By(a.validateTime))),
		validation.Field(&a.Environment, validation.Required, validation.In(
			"MIGRATION",
			"DEVELOPMENT",
			"TESTING",
			"STAGING",
			"PRODUCTION",
		)),
	)
}

func (a appConfig) validateTime(value any) error {
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

// parse method    parses a string value from env to
// the dedication type destination.
func (a *appConfig) parse() {
	raterLimit, err := strconv.Atoi(os.Getenv("RATER_LIMIT"))
	if err != nil {
		log.Fatalf("%s", err)
	}

	burstLimit, err := strconv.Atoi(os.Getenv("BURST_LIMIT"))
	if err != nil {
		log.Fatalf("%s", err)
	}

	a.RaterLimit = raterLimit
	a.BurstLimit = burstLimit

	if a.raterEvIntStr != "" {
		a.RaterEvaluationInterval = a.parseTime(a.raterEvIntStr)
	}

	if a.raterDelTime != "" {
		a.RaterDeletionTime = a.parseTime(a.raterDelTime)
	}

	defaultPaginationSize, err := strconv.Atoi(a.defaultPaginationSizeStr)
	if err != nil {
		log.Fatalf("%s", err)
	}

	a.DefaultPaginationSize = defaultPaginationSize
}

func (a appConfig) parseTime(str string) time.Duration {
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
