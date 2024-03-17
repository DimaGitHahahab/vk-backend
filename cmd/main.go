package main

import (
	"context"
	"errors"
	"fmt"
	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vk-backend/internal/api/server"
	"vk-backend/internal/repository"
	"vk-backend/internal/service/actor"
	"vk-backend/internal/service/movie"
)

func main() {
	log := logger.New()
	log.SetLevel(logger.InfoLevel)
	log.SetFormatter(&logger.TextFormatter{})

	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return fmt.Errorf("sigQuit context done")
		}
	})
	config, err := pgxpool.ParseConfig("postgres://postgres:password@localhost:5432?sslmode=disable")
	if err != nil {
		logger.Fatalf("failed to parse pgxpool config: %v", err)
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(log),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatalf("failed to create new pool: %v", err)
	}
	defer pool.Close()

	actRepo := repository.NewActorRepository(pool, log)
	movieRepo := repository.NewMovieRepository(pool, log)

	actSrv := actor.NewService(actRepo)
	movieSrv := movie.NewService(movieRepo)

	srv := server.New(":8080", &actSrv, &movieSrv, log)
	go func() {
		log.Println("starting server...")
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	if err := eg.Wait(); err != nil {
		log.Infof("gracefully shutting down the server: %v", err)
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown the server gracefully: %v", err)
	}
	log.Info("server shutdown is successful")
}
