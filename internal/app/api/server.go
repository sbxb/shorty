package api

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	srv             *http.Server
	idleConnsClosed chan struct{}
	shutdownTimeout time.Duration
	closer          func()
}

// NewHTTPServer creates a new server
func NewHTTPServer(address string, router http.Handler, closer func()) (*HTTPServer, error) {
	// Set more reasonable timeouts than the default ones
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 8 * time.Second,
		IdleTimeout:  36 * time.Second,
	}

	return &HTTPServer{
		srv:             server,
		idleConnsClosed: make(chan struct{}), // channel is closed after shutdown completed
		shutdownTimeout: 3 * time.Second,
		closer:          closer,
	}, nil
}

// WaitForInterrupt is a monitoring goroutine for catching interrupts and initiating
// server shutdown
func (s *HTTPServer) WaitForInterrupt() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-ctx.Done()
	s.Close()
}

// Close gracefully stops the server. Any additional on-close actions should be added
// here and called before idleConnsClosed channel is closed
func (s *HTTPServer) Close() {
	log.Println("Trying to gracefully stop HTTPServer")
	// Perform server shutdown with a default maximum timeout of 3 seconds
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	s.closer()
	if err := s.srv.Shutdown(timeoutCtx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTPServer Shutdown(): %v", err)
	}

	close(s.idleConnsClosed)
}

// Run starts the server
// Returns exit status code
func (s *HTTPServer) Run() int {
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTPServer ListenAndServe(): %v", err)
		return 1
	}

	<-s.idleConnsClosed
	log.Println("HTTPServer gracefully stopped")

	return 0
}
