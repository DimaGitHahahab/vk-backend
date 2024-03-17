package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository/queries"
)

type ActorRepository interface {
	AddActor(ctx context.Context, name string, gender int, birthDate time.Time) (*domain.Actor, error)
	GetActorById(ctx context.Context, id int) (*domain.Actor, error)

	ListActors(ctx context.Context) ([]*domain.Actor, error)
	UpdateActor(ctx context.Context, new *domain.Actor) error
	DeleteActor(ctx context.Context, id int) error

	ActorExists(ctx context.Context, id int) (bool, error)
}

type actorRepo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger logrus.FieldLogger
}

func NewActorRepository(pool *pgxpool.Pool, logger logrus.FieldLogger) ActorRepository {
	return &actorRepo{
		Queries: queries.NewQueries(pool),
		pool:    pool,
		logger:  logger,
	}
}
