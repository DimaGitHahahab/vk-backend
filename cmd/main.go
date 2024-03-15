package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository"
)

func main() {

	ctx := context.Background()
	logger := logrus.New()
	pool, err := pgxpool.New(ctx, "postgres://postgres:password@localhost:5432?sslmode=disable")
	if err != nil {
		logger.Fatalf("failed to create connection pool: %v", err)
	}
	movies := repository.NewMovieRepository(pool, logger)
	actors := repository.NewActorRepository(pool, logger)

	act, err := actors.AddActor(ctx, "Thomas Shelby", 1, time.Now())
	if err != nil {
		logger.Fatalf("failed to add actor: %v", err)
	}
	fmt.Println(act)

	movie, err := movies.AddMovie(ctx, "Vk-movie-2024", "Some description.", time.Now(), 9.2, []*domain.Actor{act})
	if err != nil {
		logger.Fatalf("failed to add movie: %v", err)
	}
	fmt.Println(movie)
	fmt.Println(movie.Actors[0])

	if actors.UpdateActor(ctx, &domain.Actor{
		Id:        1,
		Name:      act.Name,
		Gender:    act.Gender,
		BirthDate: act.BirthDate.AddDate(-50, 0, 0),
	}) != nil {
		logger.Fatalf("failed to update actor: %v", err)
	}
	act, err = actors.GetActorById(ctx, 1)
	if err != nil {
		logger.Fatalf("failed to get actor: %v", err)
	}
	fmt.Println(act)

	movie, err = movies.GetMovieById(ctx, 1)
	if err != nil {
		logger.Fatalf("failed to get movie: %v", err)
	}
	fmt.Println(movie)
	fmt.Println(movie.Actors[0])

}
