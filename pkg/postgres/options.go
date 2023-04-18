package postgres

import "time"

// Option -.
type Option func(*Postgres)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		if size > 0 {
			c.maxIdleConn = size
		}
	}
}

// ConnAttempts -.
func MaxOpenCoon(size int) Option {
	return func(c *Postgres) {
		if size > 0 {
			c.maxOpenConn = size
		}
	}
}

// ConnTimeout -.
func MaxConnLifetime(duration time.Duration) Option {
	return func(c *Postgres) {
		if duration > time.Duration(0) {
			c.maxLifeTime = duration
		}
	}
}
