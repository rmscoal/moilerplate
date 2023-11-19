package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type dbConfig struct {
	URL      string
	host     string
	port     string
	user     string
	password string
	name     string

	maxPoolSize     int
	maxOpenConn     int
	maxConnLifetime time.Duration
}

// newDbConfig method    has a receiver of the config
// struct. It loads the dbConfig struct into the main
// Config struct.
func (c *Config) newDbConfig() {
	d := dbConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		name:     os.Getenv("DB_NAME"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
	}

	if err := (&d).parse(); err != nil {
		log.Fatalf("Error parsing postgres environment: %s\n", err)
	}

	if err := d.validate(); err != nil {
		log.Fatalf("%s", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		d.user,
		d.password,
		d.host,
		d.port,
		d.name,
	)
	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatalf("ERROR parsing dsn: %s\n", err)
	}
	u.User = url.UserPassword(d.user, d.password)

	c.Db.URL = u.String()
	c.Db = d
}

// validate method    validates the dbConfig struct
// such that in matches the expected conditions.
func (d dbConfig) validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.host, validation.Required, is.Host),
		validation.Field(&d.port, validation.Required, is.Port.Error("hello world")),
		validation.Field(&d.user, validation.Required),
		validation.Field(&d.name, validation.Required),
		validation.Field(&d.password, validation.Required),
	)
}

// parse method    parses a string value from env to
// the dedication type destination.
func (d *dbConfig) parse() (err error) {
	d.maxPoolSize, err = strconv.Atoi(os.Getenv("DB_MAX_POOL_SIZE"))
	if err != nil {
		return err
	}

	d.maxOpenConn, err = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))
	if err != nil {
		return err
	}

	d.maxConnLifetime, err = time.ParseDuration(os.Getenv("DB_MAX_CONN_LIFETIME"))
	if err != nil {
		return err
	}

	return nil
}

// Getter functions for getting db connection configurations
func (d dbConfig) MaxPoolSize() int {
	return d.maxPoolSize
}

func (d dbConfig) MaxOpenConn() int {
	return d.maxOpenConn
}

func (d dbConfig) MaxConnLifetime() time.Duration {
	return d.maxConnLifetime
}
