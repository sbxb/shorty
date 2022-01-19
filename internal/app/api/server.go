package api

import (
	"context"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	srv             *http.Server
	idleConnsClosed chan struct{}
	shutdownTimeout time.Duration
}

// NewHTTPServer creates a new server
func NewHTTPServer(address string, router http.Handler) (*HTTPServer, error) {
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
	}, nil
}

// Close gracefully stops the server. Any additional on-close actions should be added
// here and called before idleConnsClosed channel is closed
func (s *HTTPServer) Close() {
	if s.srv == nil {
		return
	}
	log.Println("Trying to gracefully stop HTTPServer")
	// Perform server shutdown with a default maximum timeout of 3 seconds
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(timeoutCtx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTPServer Shutdown() failed: %v", err)
	}

	s.srv = nil
	close(s.idleConnsClosed)
}

// Start runs the server and creates a monitoring gorouting to wait for
// the context to be marked done
func (s *HTTPServer) Start(ctx context.Context) {
	if s.srv == nil {
		return
	}
	go func() {
		<-ctx.Done()
		s.Close()
	}()

	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTPServer ListenAndServe() failed: %v", err)
		return
	}

	<-s.idleConnsClosed
	log.Println("HTTPServer gracefully stopped")
}
