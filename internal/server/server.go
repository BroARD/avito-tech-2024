package server

import (
	"avito-tech/config"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg *config.Config
	db *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		cfg: cfg,
		db: db,
	}
}

func (s *Server) Run() error {
	router := s.MapRoutes()

	httpServer := &http.Server{
		Addr: ":" + s.cfg.ServPort,
		Handler: router,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalln("GG Server UPAL")
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при graceful shutdown: %v", err)
		return err
	}

	return nil
}
