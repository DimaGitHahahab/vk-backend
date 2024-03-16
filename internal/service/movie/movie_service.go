package movie

import (
	"context"
	"fmt"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository"
	"vk-backend/internal/service/util"
)

type Service interface {
	AddMovie(ctx context.Context, title string, description string, releaseDate time.Time, rating float64, actors []*domain.Actor) (*domain.Movie, error)
	GetMovieById(ctx context.Context, id int) (*domain.Movie, error)
	GetActorsByMovieId(ctx context.Context, movieId int) ([]*domain.Actor, error)
	ListMovies(ctx context.Context, filter *util.Filter, sorting util.SortBy) ([]*domain.Movie, error)
	UpdateMovie(ctx context.Context, new *domain.Movie) error
	DeleteMovie(ctx context.Context, id int) error
}

type movieService struct {
	repo repository.MovieRepository
}

func NewService(repo repository.MovieRepository) Service {
	return &movieService{
		repo: repo,
	}
}

func (s *movieService) AddMovie(ctx context.Context, title string, description string, releaseDate time.Time, rating float64, actors []*domain.Actor) (*domain.Movie, error) {
	err := util.ValidateMovieData(title, description, rating)
	if err != nil {
		return nil, err
	}
	movie, err := s.repo.AddMovie(ctx, title, description, releaseDate, rating, actors)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *movieService) GetMovieById(ctx context.Context, id int) (*domain.Movie, error) {
	ok, err := s.repo.MovieExists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("movie service can't check if movie exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrMovieNotExists
	}

	movie, err := s.repo.GetMovieById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("movie service can't get movie by id: %w", err)
	}

	return movie, nil
}

func (s *movieService) GetActorsByMovieId(ctx context.Context, movieId int) ([]*domain.Actor, error) {
	ok, err := s.repo.MovieExists(ctx, movieId)
	if err != nil {
		return nil, fmt.Errorf("movie service can't check if movie exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrMovieNotExists
	}

	actors, err := s.repo.GetActorsByMovieId(ctx, movieId)
	if err != nil {
		return nil, fmt.Errorf("movie service can't get actors by movie id: %w", err)
	}

	return actors, nil
}

func (s *movieService) ListMovies(ctx context.Context, filter *util.Filter, sorting util.SortBy) ([]*domain.Movie, error) {
	movies, err := s.repo.ListMovies(ctx)
	if err != nil {
		return nil, fmt.Errorf("movie service can't list movies: %w", err)
	}
	movies = util.SortMovies(util.FilterMovies(movies, filter), sorting)

	return movies, nil
}

func (s *movieService) UpdateMovie(ctx context.Context, new *domain.Movie) error {
	ok, err := s.repo.MovieExists(ctx, new.Id)
	if err != nil {
		return fmt.Errorf("movie service can't check if movie exists: %w", err)
	}
	if !ok {
		return domain.ErrMovieNotExists
	}

	err = s.repo.UpdateMovie(ctx, new)
	if err != nil {
		return fmt.Errorf("movie service can't update movie: %w", err)
	}

	return nil
}

func (s *movieService) DeleteMovie(ctx context.Context, id int) error {
	ok, err := s.repo.MovieExists(ctx, id)
	if err != nil {
		return fmt.Errorf("movie service can't check if movie exists: %w", err)
	}
	if !ok {
		return domain.ErrMovieNotExists
	}

	err = s.repo.DeleteMovie(ctx, id)
	if err != nil {
		return fmt.Errorf("movie service can't delete movie: %w", err)
	}

	return nil
}
