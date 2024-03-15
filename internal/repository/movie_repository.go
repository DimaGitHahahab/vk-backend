package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository/queries"
)

type MovieRepository interface {
	AddMovie(ctx context.Context, title string, description string, releaseDate time.Time, rating float64, actors []*domain.Actor) (*domain.Movie, error)
	GetMovieById(ctx context.Context, id int) (*domain.Movie, error)
	ListMovies(ctx context.Context) ([]*domain.Movie, error)
	FindMoviesByTitle(ctx context.Context, name string) ([]*domain.Movie, error)
	FindMoviesByActorName(ctx context.Context, name string) ([]*domain.Movie, error)
	GetMoviesWithActor(ctx context.Context, actorId int) ([]*domain.Movie, error)
	UpdateMovie(ctx context.Context, new *domain.Movie) error
	DeleteMovie(ctx context.Context, id int) error
}

type movieRepo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger logrus.FieldLogger
}

func NewMovieRepository(pool *pgxpool.Pool, logger logrus.FieldLogger) MovieRepository {
	return &movieRepo{
		Queries: queries.NewQueries(pool),
		pool:    pool,
		logger:  logger,
	}
}
