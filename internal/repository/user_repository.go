package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository/queries"
)

type UserRepository interface {
	AddUser(ctx context.Context, name string, password string) (*domain.User, error)
	UserExists(ctx context.Context, name string) (bool, error)
	GetUserByName(ctx context.Context, name string) (*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
}

type userRepo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger logrus.FieldLogger
}

func NewUserRepository(pool *pgxpool.Pool, logger logrus.FieldLogger) UserRepository {
	return &userRepo{
		Queries: queries.NewQueries(pool),
		pool:    pool,
		logger:  logger,
	}
}
