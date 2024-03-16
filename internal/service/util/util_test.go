package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"vk-backend/internal/domain"
)

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
