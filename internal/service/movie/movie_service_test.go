package movie

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"strings"
	"testing"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/mocks"
)

func TestMovieService_AddMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	releaseDate := time.Now()
	repo.
		EXPECT().
		AddMovie(gomock.Any(), "title", "description", releaseDate, 9.0, nil).
		Return(&domain.Movie{
			Id:          1,
			Title:       "title",
			Description: "description",
			ReleaseDate: releaseDate,
			Rating:      9.0,
			Actors:      nil,
		}, nil)

	movie, err := service.AddMovie(context.Background(), "title", "description", releaseDate, 9.0, nil)
	assert.NoError(t, err)
	assert.Equal(t, &domain.Movie{
		Id:          1,
		Title:       "title",
		Description: "description",
		ReleaseDate: releaseDate,
		Rating:      9.0,
		Actors:      nil,
	}, movie)
}

func TestMovieService_AddMovie_WithActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	releaseDate := time.Now()
	actors := []*domain.Actor{
		{
			Id:        1,
			Name:      "name",
			Gender:    1,
			BirthDate: time.Now(),
		},
	}
	repo.
		EXPECT().
		AddMovie(gomock.Any(), "title", "description", releaseDate, 9.0, actors).
		Return(&domain.Movie{
			Id:          1,
			Title:       "title",
			Description: "description",
			ReleaseDate: releaseDate,
			Rating:      9.0,
			Actors:      actors,
		}, nil)

	movie, err := service.AddMovie(context.Background(), "title", "description", releaseDate, 9.0, actors)
	assert.NoError(t, err)
	assert.Equal(t, &domain.Movie{
		Id:          1,
		Title:       "title",
		Description: "description",
		ReleaseDate: releaseDate,
		Rating:      9.0,
		Actors:      actors,
	}, movie)
}

func TestMovieService_AddMovie_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	releaseDate := time.Now()
	movie, err := service.AddMovie(context.Background(), "", "description", releaseDate, 9.0, nil)
	assert.ErrorIs(t, err, domain.ErrEmptyTitle)
	assert.Nil(t, movie)

	movie, err = service.AddMovie(context.Background(), "title", "", releaseDate, 9.0, nil)
	assert.ErrorIs(t, err, domain.ErrEmptyDescription)
	assert.Nil(t, movie)

	movie, err = service.AddMovie(context.Background(), "title", "description", releaseDate, -2.0, nil)
	assert.ErrorIs(t, err, domain.ErrInvalidRating)
	assert.Nil(t, movie)

	longTitle := strings.Repeat("a", 256)
	movie, err = service.AddMovie(context.Background(), longTitle, "description", releaseDate, 9.0, nil)
	assert.ErrorIs(t, err, domain.ErrTooLongTitle)
	assert.Nil(t, movie)

	longDescription := strings.Repeat("a", 4096)
	movie, err = service.AddMovie(context.Background(), "title", longDescription, releaseDate, 9.0, nil)
	assert.ErrorIs(t, err, domain.ErrTooLongDescription)
	assert.Nil(t, movie)
}

func TestMovieService_GetMovieById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	releaseDate := time.Now()
	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(true, nil)
	repo.
		EXPECT().
		GetMovieById(gomock.Any(), 1).
		Return(&domain.Movie{
			Id:          1,
			Title:       "title",
			Description: "description",
			ReleaseDate: releaseDate,
			Rating:      9.0,
			Actors:      nil,
		}, nil)

	movie, err := service.GetMovieById(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, &domain.Movie{
		Id:          1,
		Title:       "title",
		Description: "description",
		ReleaseDate: releaseDate,
		Rating:      9.0,
		Actors:      nil,
	}, movie)
}

func TestMovieService_GetMovieById_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(false, nil)

	movie, err := service.GetMovieById(context.Background(), 1)
	assert.ErrorIs(t, err, domain.ErrMovieNotExists)
	assert.Nil(t, movie)
}

func TestMovieService_GetActorsByMovieId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	birthDate := time.Now()
	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(true, nil)
	repo.
		EXPECT().
		GetActorsByMovieId(gomock.Any(), 1).
		Return([]*domain.Actor{
			{
				Id:        1,
				Name:      "name",
				Gender:    1,
				BirthDate: birthDate,
			},
			{
				Id:        2,
				Name:      "name2",
				Gender:    2,
				BirthDate: birthDate.AddDate(1, 0, 0),
			},
		}, nil)

	actors, err := service.GetActorsByMovieId(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, []*domain.Actor{
		{
			Id:        1,
			Name:      "name",
			Gender:    1,
			BirthDate: birthDate,
		},
		{
			Id:        2,
			Name:      "name2",
			Gender:    2,
			BirthDate: birthDate.AddDate(1, 0, 0),
		},
	}, actors)
}

func TestMovieService_GetActorsByMovieId_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(false, nil)

	actors, err := service.GetActorsByMovieId(context.Background(), 1)
	assert.ErrorIs(t, err, domain.ErrMovieNotExists)
	assert.Nil(t, actors)
}

func TestMovieService_UpdateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	releaseDate := time.Now()

	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(true, nil)

	repo.
		EXPECT().
		UpdateMovie(gomock.Any(), &domain.Movie{
			Id:          1,
			Title:       "title",
			Description: "description",
			ReleaseDate: releaseDate,
			Rating:      9.0,
			Actors:      nil,
		}).
		Return(nil)

	err := service.UpdateMovie(context.Background(), &domain.Movie{
		Id:          1,
		Title:       "title",
		Description: "description",
		ReleaseDate: releaseDate,
		Rating:      9.0,
		Actors:      nil,
	})

	assert.NoError(t, err)
}

func TestMovieService_UpdateMovie_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	releaseDate := time.Now()

	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(false, nil)

	err := service.UpdateMovie(context.Background(), &domain.Movie{
		Id:          1,
		Title:       "title",
		Description: "description",
		ReleaseDate: releaseDate,
		Rating:      9.0,
		Actors:      nil,
	})

	assert.ErrorIs(t, err, domain.ErrMovieNotExists)
}

func TestMovieService_DeleteMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(true, nil)

	repo.
		EXPECT().
		DeleteMovie(gomock.Any(), 1).
		Return(nil)

	err := service.DeleteMovie(context.Background(), 1)
	assert.NoError(t, err)
}

func TestMovieService_DeleteMovie_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		MovieExists(gomock.Any(), 1).
		Return(false, nil)

	err := service.DeleteMovie(context.Background(), 1)
	assert.ErrorIs(t, err, domain.ErrMovieNotExists)
}