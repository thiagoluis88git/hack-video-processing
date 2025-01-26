package httpserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	_defaultReadTimeout     = 20 * time.Second
	_defaultWriteTimeout    = 20 * time.Second
	_defaultShutdownTimeout = 6 * time.Second
	port                    = 3211
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         fmt.Sprintf(":%d", port),
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	return s
}

func (s *Server) Start() {
	go func() {
		listener, err := net.Listen("tcp", s.server.Addr)

		if err != nil {
			//shutdown
			err := s.Shutdown()
			if err != nil {
				log.Print("httpServer shutdown", map[string]interface{}{
					"signal": err.Error(),
				})
			}
		}

		log.Print("Video Processing API has started")

		s.notify <- s.server.Serve(listener)
		close(s.notify)
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case signalInterrupt := <-interrupt:
		log.Print("signal interrupt received", map[string]interface{}{
			"signal": signalInterrupt.String(),
		})
	case err := <-s.Notify():
		log.Print("httpServer notify and error", map[string]interface{}{
			"notify": err.Error(),
		})
	}

	// Shutdown
	err := s.Shutdown()
	if err != nil {
		log.Print("httpServer shutdown", map[string]interface{}{
			"notify": err.Error(),
		})
	}
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
