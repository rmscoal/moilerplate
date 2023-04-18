package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":8081"
	_defaultShutdownTimeout = 3 * time.Second
)

var (
	once                 sync.Once
	serverSingleInstance *Server
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
	ran             bool
}

// NewServer function    creates a new server and
// calls `run` internally which starts the server
// and handles graceful shutdown automatically.
func NewServer(handler http.Handler, opts ...Option) *Server {
	if serverSingleInstance == nil {
		once.Do(func() {
			httpServer := &http.Server{
				Handler:      handler,
				ReadTimeout:  _defaultReadTimeout,
				WriteTimeout: _defaultWriteTimeout,
				Addr:         _defaultAddr,
			}

			serverSingleInstance = &Server{
				server:          httpServer,
				shutdownTimeout: _defaultShutdownTimeout,
			}

			for _, opt := range opts {
				opt(serverSingleInstance)
			}

			serverSingleInstance.run()
		})
	}

	return serverSingleInstance
}

func (s *Server) CheckConn() error {
	if !s.ran {
		if serverSingleInstance == nil {
			return errors.New("Server has not yet been initialized")
		}
	}

	return nil
}

// run method    runs the server and also handles
// the graceful shutdown.
func (s *Server) run() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}

		s.ran = true
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server ... Please wait ⌛️")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("FATAL - Error while shutting down server: %s", err)
	} else {
		log.Println("INFO - Server successfully shutdown")
		cancel()
	}

	<-ctx.Done()
	log.Println("Server closed")

	log.Println("Exiting")
}
