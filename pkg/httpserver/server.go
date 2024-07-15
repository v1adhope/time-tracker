package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Socket          string `koanf:"APP_SERVER_SOCKET"`
	ShutdownTimeout int64  `koanf:"APP_SERVER_SHUTDOWN_TIMEOUT"`
	ReadTimeout     int64  `koanf:"APP_SERVER_READ_TIMEOUT"`
	WriteTimeout    int64  `koanf:"APP_SERVER_WRITE_TIMEOUT"`
}

type Server struct {
	server          *http.Server
	shutdownTimeout int64
}

func parseDuration(t int64) time.Duration {
	return time.Duration(t) * time.Second
}

func New(handler http.Handler, cfg *Config) *Server {
	httpServer := &http.Server{
		Addr:         cfg.Socket,
		Handler:      handler,
		ReadTimeout:  parseDuration(cfg.ReadTimeout),
		WriteTimeout: parseDuration(cfg.WriteTimeout),
	}

	return &Server{
		server:          httpServer,
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			log.Fatalf("httpserver: listenAndServe: %s", err)
		}
	}()

	s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), parseDuration(s.shutdownTimeout))
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("httpserver: gracefulShutdown: shutdown: %s", err)
	}

	select {
	case <-ctx.Done():
		log.Printf("timeout of %d seconds", s.shutdownTimeout)
	}

	log.Print("server exiting")
}
