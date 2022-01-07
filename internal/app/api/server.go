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
}

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
		shutdownTimeout: 1 * time.Second,
	}, nil
}

func (s *HTTPServer) WaitForInterrupt() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-ctx.Done()
	log.Println("Interrupt caught")
	s.Close()
}

func (s *HTTPServer) Close() {
	log.Println("Trying to gracefully stop HTTPServer")
	// Perform server shutdown with a default maximum timeout of 1 second
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(timeoutCtx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTPServer Shutdown(): %v", err)
	}
	close(s.idleConnsClosed)
	log.Println("HTTPServer gracefully stopped")
}

func (s *HTTPServer) Run() int {
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTPServer ListenAndServe(): %v", err)
		return 1
	}
	<-s.idleConnsClosed
	return 0
}
