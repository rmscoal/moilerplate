package httpserver

import (
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
func RegisterHostAndPort(host, port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort(host, port)
	}
}

// ReadTimeout -.
func RegisterReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout -.
func RegisterWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// ShutdownTimeout -.
func RegisterShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
