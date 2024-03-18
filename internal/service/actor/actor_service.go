package actor

import (
	"context"
	"fmt"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository"
)

type Service interface {
	AddActor(ctx context.Context, name string, gender int, birthDate time.Time) (*domain.Actor, error)
	GetActorById(ctx context.Context, id int) (*domain.Actor, error)

	UpdateActor(ctx context.Context, new *domain.Actor) error
	DeleteActor(ctx context.Context, id int) error

	ListActors(ctx context.Context) ([]*domain.Actor, error)
}

type actorService struct {
	repo repository.ActorRepository
}

func NewService(repo repository.ActorRepository) Service {
	return &actorService{
		repo: repo,
	}
}
func (s *actorService) AddActor(ctx context.Context, name string, gender int, birthDate time.Time) (*domain.Actor, error) {
	err := validateActorData(name, birthDate, gender)
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
	if id <= 0 {
		return nil, domain.ErrActorNotExists
	}
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

func (s *actorService) UpdateActor(ctx context.Context, new *domain.Actor) error {
	if new.Id <= 0 {
		return domain.ErrActorNotExists
	}

	ok, err := s.repo.ActorExists(ctx, new.Id)
	if err != nil {
		return fmt.Errorf("actor service can't check if actor exists: %w", err)
	}
	if !ok {
		return domain.ErrActorNotExists
	}

	err = validateActorData(new.Name, new.BirthDate, new.Gender)
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
	if id <= 0 {
		return domain.ErrActorNotExists
	}
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

func (s *actorService) ListActors(ctx context.Context) ([]*domain.Actor, error) {
	actors, err := s.repo.ListActors(ctx)
	if err != nil {
		return nil, fmt.Errorf("actor service can't list actors: %w", err)
	}

	return actors, nil

}

func validateActorData(name string, birthDate time.Time, gender int) error {
	if name == "" {
		return domain.ErrEmptyName
	}
	if birthDate.After(time.Now()) {
		return domain.ErrFutureBirthDate
	}
	if birthDate.IsZero() {
		return domain.ErrEmptyBirthDate
	}
	if gender != 0 && gender != 1 && gender != 2 && gender != 9 {
		return domain.ErrInvalidGender
	}

	return nil
}
