package service

import (
	"context"
	"fmt"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository"
	"vk-backend/internal/service/util"
)

type ActorService interface {
	AddActor(ctx context.Context, name string, gender int, birthDate time.Time) (*domain.Actor, error)
	GetActorById(ctx context.Context, id int) (*domain.Actor, error)

	AddActorToMovie(ctx context.Context, actorId int, movieId int) error

	UpdateActor(ctx context.Context, new *domain.Actor) error
	DeleteActor(ctx context.Context, id int) error
}

type actorService struct {
	repo repository.ActorRepository
}

func NewActorService(repo repository.ActorRepository) ActorService {
	return &actorService{
		repo: repo,
	}
}
func (s *actorService) AddActor(ctx context.Context, name string, gender int, birthDate time.Time) (*domain.Actor, error) {
	err := util.ValidateActorData(name, birthDate)
	if err != nil {
		return nil, err
	}
	actor, err := s.repo.AddActor(ctx, name, gender, birthDate)
	if err != nil {
		return nil, fmt.Errorf("actor service can't add actor: %w", err)
	}

	return actor, nil
}

func (s *actorService) GetActorById(ctx context.Context, id int) (*domain.Actor, error) {
	ok, err := s.repo.ActorExists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("actor service can't check if actor exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrActorNotExists
	}

	actor, err := s.repo.GetActorById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("actor service can't get actor by id: %w", err)
	}

	return actor, nil
}
func (s *actorService) AddActorToMovie(ctx context.Context, actorId int, movieId int) error {
	ok, err := s.repo.ActorExists(ctx, actorId)
	if err != nil {
		return fmt.Errorf("actor service can't check if actor exists: %w", err)
	}
	if !ok {
		return domain.ErrActorNotExists
	}

	ok, err = s.repo.MovieExists(ctx, movieId)
	if err != nil {
		return fmt.Errorf("actor service can't check if movie exists: %w", err)
	}
	if !ok {
		return domain.ErrMovieNotExists
	}

	err = s.repo.AddActorToMovie(ctx, actorId, movieId)
	if err != nil {
		return fmt.Errorf("actor service can't add actor to movie: %w", err)
	}

	return nil
}

func (s *actorService) UpdateActor(ctx context.Context, new *domain.Actor) error {
	ok, err := s.repo.ActorExists(ctx, new.Id)
	if err != nil {
		return fmt.Errorf("actor service can't check if actor exists: %w", err)
	}
	if !ok {
		return domain.ErrActorNotExists
	}

	err = util.ValidateActorData(new.Name, new.BirthDate)
	if err != nil {
		return err
	}

	err = s.repo.UpdateActor(ctx, new)
	if err != nil {
		return fmt.Errorf("actor service can't update actor: %w", err)
	}

	return nil
}

func (s *actorService) DeleteActor(ctx context.Context, id int) error {
	ok, err := s.repo.ActorExists(ctx, id)
	if err != nil {
		return fmt.Errorf("actor service can't check if actor exists: %w", err)
	}
	if !ok {
		return domain.ErrActorNotExists
	}

	err = s.repo.DeleteActor(ctx, id)
	if err != nil {
		return fmt.Errorf("actor service can't delete actor: %w", err)
	}

	return nil
}
