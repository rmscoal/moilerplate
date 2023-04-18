// Postgres is a part of the database. This will generate a new
// connection with PostgreSQL database. By default it construct
// an *sql.DB and passes the connection pool to GORM for later
// use. By doing this, it enables developer to use any ORM or
// tools to query since we are still using the built-in *sql.DB.
// To pass in database config from the environment, do explicitly
// pass optionals to the parameter on creating the Postgres
// instance.
//
// We also follow the singleton creation design pattern since
// sql.Open will create a new pool.

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once             sync.Once
	pgSingleInstance *Postgres
)

var (
	_defaultMaxPoolSize  = 2
	_defaultConnAttempts = 10
	_defautlMaxOpenConn  = 10
	_defaultConnLifeTime = 5 * time.Minute
)

// The pool size can be controlled with SetMaxIdleConns.
type Postgres struct {
	maxIdleConn     int           // manages the pool size for a connections
	maxOpenConn     int           // manages the maximum number of connectios to the database
	maxLifeTime     time.Duration // manages the maximum amount of time connection may be reused.
	maxConnAttempts int           // manages the maximum attemp to ping the database using the connection

	ORM  *gorm.DB
	Pool *sql.DB

	initialized bool
}

// NewPostgres function    returns a new Postgres with
// a custom option to set the *sql.DB configurations.
// It has to be explicitly passed into the parameter
// while calling this function.

// Change to *config.Config
func GetPostgres(url string, opts ...Option) *Postgres {
	if pgSingleInstance == nil {
		log.Println("Creating postgres instance")
		once.Do(func() {
			pgSingleInstance = &Postgres{
				maxIdleConn:     _defaultMaxPoolSize,
				maxOpenConn:     _defautlMaxOpenConn,
				maxLifeTime:     _defaultConnLifeTime,
				maxConnAttempts: _defaultConnAttempts,
				initialized:     false,
			}

			for _, opt := range opts {
				opt(pgSingleInstance)
			}

			pgSingleInstance.connect(url)
			pgSingleInstance.ping()
		})
	}

	return pgSingleInstance
}

func (pg *Postgres) CheckConn() error {
	if !pg.initialized {
		if pgSingleInstance == nil {
			return errors.New("Postgres has not yet been initialized")
		}
	}

	return nil
}

func (pg *Postgres) connect(url string) {
	sqldb, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("FATAL - Unable to start database")
	}

	sqldb.SetMaxIdleConns(pg.maxIdleConn)
	sqldb.SetMaxOpenConns(pg.maxOpenConn)
	sqldb.SetConnMaxLifetime(pg.maxLifeTime)

	pg.Pool = sqldb

	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL - Unable to connect to postgres: %s", err)
	}

	pg.ORM = gormdb
	pg.initialized = true
}

func (pg *Postgres) ping() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for pg.maxConnAttempts > 0 {
		pg.maxConnAttempts--

		err := pg.Pool.PingContext(ctx)
		if err == nil {
			log.Printf("INFO - Successfully connected to postgreSQL database after %d attempt(s)", _defaultConnAttempts-pg.maxConnAttempts)
			return
		}
	}

	log.Fatalf("FATAL - Unable to ping database")
}

func (pg *Postgres) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
