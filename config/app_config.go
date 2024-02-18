package config

import (
	"log"
	"os"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type appConfig struct {
	raterLimit              int
	burstLimit              int
	raterEvaluationInterval time.Duration
	raterDeletionTime       time.Duration
	defaultPaginationSize   int
}

// newServerConfig method    has a Config receiver
// such that it loads the serverConfig to the main
// Config.
func (c *Config) newAppConfig() {
	var a appConfig

	if err := a.parse(); err != nil {
		log.Fatalf("Error while parsing app configuration: %s\n", err)
	}

	if err := a.validate(); err != nil {
		log.Fatalf("FATAL - %s", err)
	}

	c.App = a
}

// validate method    validates the serverConfig
// values such that it meets the requirements.
func (a appConfig) validate() error {
	return validation.ValidateStruct(&a)
}

// parse method    parses a string value from env to
// the dedication type destination.
func (a *appConfig) parse() (err error) {
	a.raterLimit, err = strconv.Atoi(os.Getenv("RATER_LIMIT"))
	if err != nil {
		return err
	}

	a.burstLimit, err = strconv.Atoi(os.Getenv("BURST_LIMIT"))
	if err != nil {
		return err
	}

	a.raterEvaluationInterval, err = time.ParseDuration(os.Getenv("RATER_EVALUATION_INTERVAL"))
	if err != nil {
		return err
	}

	a.raterDeletionTime, err = time.ParseDuration(os.Getenv("RATER_DELETION_TIME"))
	if err != nil {
		return err
	}

	a.defaultPaginationSize, err = strconv.Atoi(os.Getenv("DEFAULT_ROWS_PER_PAGE"))
	if err != nil {
		return err
	}

	return nil
}

// Getter functions for getting app configurations

func (a appConfig) RaterLimit() int {
	return a.raterLimit
}

func (a appConfig) BurstLimit() int {
	return a.burstLimit
}

func (a appConfig) DefaultPaginationSize() int {
	return a.defaultPaginationSize
}

func (a appConfig) RaterEvaluationInterval() time.Duration {
	return a.raterEvaluationInterval
}

func (a appConfig) RaterDeletionTime() time.Duration {
	return a.raterDeletionTime
}
