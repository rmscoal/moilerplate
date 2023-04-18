package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type dbConfig struct {
	URL      string
	host     string
	port     string
	user     string
	password string
	name     string

	maxPoolSizeStr     string
	maxOpenConnStr     string
	maxConnLifetimeStr string

	maxPoolSize     int
	maxOpenConn     int
	maxConnLifetime int
}

// newDbConfig method    has a receiver of the config
// struct. It loads the dbConfig struct into the main
// Config struct.
func (c *Config) newDbConfig() {
	d := dbConfig{
		host:               os.Getenv("DB_HOST"),
		port:               os.Getenv("DB_PORT"),
		name:               os.Getenv("DB_NAME"),
		user:               os.Getenv("DB_USER"),
		password:           os.Getenv("DB_PASSWORD"),
		maxPoolSizeStr:     os.Getenv("DB_MAX_OPEN_CONN"),
		maxOpenConnStr:     os.Getenv("DB_MAX_POOL_SIZE"),
		maxConnLifetimeStr: os.Getenv("DB_MAX_CONN_LIFETIME"),
	}

	if err := d.validate(); err != nil {
		log.Fatalf("%s", err)
	}

	pd := &d
	pd.parse()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		d.user,
		d.password,
		d.host,
		d.port,
		d.name,
	)

	c.Db = *pd
	c.Db.URL = dsn
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
		validation.Field(&d.maxOpenConnStr, is.Int),
		validation.Field(&d.maxPoolSizeStr, is.Int),
		validation.Field(&d.maxConnLifetimeStr, is.Int),
	)
}

// parse method    parses a string value from env to
// the dedication type destination.
func (d *dbConfig) parse() {
	var maxPoolSize, maxOpenConn, maxConnLifetime int

	maxPoolSize, err := strconv.Atoi(d.maxPoolSizeStr)
	if err != nil {
		maxPoolSize = -1
	}
	maxOpenConn, err = strconv.Atoi(d.maxOpenConnStr)
	if err != nil {
		maxOpenConn = -1
	}
	maxConnLifetime, err = strconv.Atoi(d.maxConnLifetimeStr)
	if err != nil {
		maxConnLifetime = -1
	}

	d.maxPoolSize = maxPoolSize
	d.maxOpenConn = maxOpenConn
	d.maxConnLifetime = maxConnLifetime
}

// Getter functions for getting db connection configurations
func (d dbConfig) MaxPoolSize() int {
	return d.maxPoolSize
}

func (d dbConfig) MaxOpenConn() int {
	return d.maxOpenConn
}

func (d dbConfig) MaxConnLifetime() int {
	return d.maxConnLifetime
}
