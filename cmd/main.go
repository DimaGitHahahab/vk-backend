package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vk-backend/internal/api/server"
	"vk-backend/internal/repository"
	"vk-backend/internal/service/actor"
	"vk-backend/internal/service/movie"
	"vk-backend/internal/service/user"
)

func main() {
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{})
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("failed to load env variables: %v", err)
	}

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

	logger.Info("starting migration...")
	processMigration(os.Getenv("MIGRATION_PATH"), os.Getenv("DB_URL"), logger)
	logger.Info("migration successful")

	config, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatalf("failed to parse pgxpool config: %v", err)
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(logger),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatalf("failed to create new pool: %v", err)
	}
	defer pool.Close()

	actRepo := repository.NewActorRepository(pool, logger)
	movieRepo := repository.NewMovieRepository(pool, logger)
	userRepo := repository.NewUserRepository(pool, logger)

	actSrv := actor.NewService(actRepo)
	movieSrv := movie.NewService(movieRepo)
	userSrv := user.NewService(userRepo)

	srv := server.New(os.Getenv("HTTP_PORT"), &actSrv, &movieSrv, &userSrv, logger)
	go func() {
		logger.Println("starting server...")
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("failed to start server: %v", err)
		}
	}()

	if err := eg.Wait(); err != nil {
		logger.Infof("gracefully shutting down the server: %v", err)
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Fatalf("failed to shutdown the server gracefully: %v", err)
	}
	logger.Info("server shutdown is successful")
}

func processMigration(migrationURL string, dbSource string, logger *log.Logger) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Fatalf("failed to create migration: %v", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatalf("failed to migrate: %v", err)
	}
	defer migration.Close()
}
