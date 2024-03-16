package movie

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"sort"
	"strings"
	"testing"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/mocks"
)

func testMovies() []*domain.Movie {
	return []*domain.Movie{
		{
			Id:          1,
			Title:       "title1",
			Description: "description1",
			ReleaseDate: time.Now(),
			Rating:      5.0,
		},
		{
			Id:          2,
			Title:       "title2",
			Description: "description2",
			ReleaseDate: time.Now().AddDate(1, 0, 0),
			Rating:      4.0,
		},
		{
			Id:          3,
			Title:       "title3",
			Description: "description3",
			ReleaseDate: time.Now().AddDate(2, 0, 0),
			Rating:      3.0,
		},
	}
}

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

func TestMovieService_ListMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	service := NewService(repo)

	movies := testMovies()

	repo.
		EXPECT().
		ListMovies(gomock.Any()).
		Return(movies, nil)

	movies, err := service.ListMovies(context.Background(), nil, DefaultSort)
	assert.NoError(t, err)
	expected := movies
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Rating > expected[j].Rating
	})
	assert.Equal(t, expected, movies)

}

func TestFilterBuilder(t *testing.T) {

	f := NewFilter()
	assert.NotNil(t, f)
	assert.Nil(t, f.title)
	assert.Nil(t, f.releaseDate)
	assert.Nil(t, f.rating)

	title, date, rating := "title", time.Now(), 5.0
	f = f.WithTitle(title).WithReleaseDate(date).WithRating(rating)

	assert.Equal(t, title, *f.title)
	assert.Equal(t, date, *f.releaseDate)
	assert.Equal(t, rating, *f.rating)
}

func TestFilterMovies(t *testing.T) {
	movies := testMovies()
	filter := NewFilter().WithTitle("title1")

	filteredMovies := FilterMovies(movies, filter)
	assert.Len(t, filteredMovies, 1)
	assert.Equal(t, "title1", filteredMovies[0].Title)

	filter = NewFilter().WithRating(4.0)
	filteredMovies = FilterMovies(movies, filter)
	assert.Len(t, filteredMovies, 2)

	filter = NewFilter().WithReleaseDate(movies[1].ReleaseDate)
	filteredMovies = FilterMovies(movies, filter)
	assert.Len(t, filteredMovies, 2)
}

func TestSortMovies(t *testing.T) {
	movies := testMovies()

	sortedMovies := SortMovies(movies, SortByRating)
	assert.Equal(t, 5.0, sortedMovies[0].Rating)
	assert.Equal(t, 4.0, sortedMovies[1].Rating)
	assert.Equal(t, 3.0, sortedMovies[2].Rating)

	sortedMovies = SortMovies(movies, SortByReleaseDate)
	assert.True(t, movies[2].ReleaseDate.Before(movies[1].ReleaseDate) && movies[1].ReleaseDate.Before(movies[0].ReleaseDate))

	sortedMovies = SortMovies(movies, SortByTitle)
	assert.Equal(t, "title1", sortedMovies[0].Title)
	assert.Equal(t, "title2", sortedMovies[1].Title)
	assert.Equal(t, "title3", sortedMovies[2].Title)
}
