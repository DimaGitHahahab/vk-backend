package movie

import (
	"strings"
	"time"
	"vk-backend/internal/domain"
)

type Filter struct {
	title       *string
	releaseDate *time.Time
	rating      *float64
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) WithTitle(title string) *Filter {
	f.title = &title
	return f
}

func (f *Filter) WithReleaseDate(releaseDate time.Time) *Filter {
	f.releaseDate = &releaseDate
	return f
}

func (f *Filter) WithRating(rating float64) *Filter {
	f.rating = &rating
	return f
}

func FilterMovies(movies []*domain.Movie, filter *Filter) []*domain.Movie {
	res := make([]*domain.Movie, 0, len(movies))
	if filter == nil {
		return movies
	}
	for _, movie := range movies {
		if filter.title != nil && !strings.Contains(movie.Title, *filter.title) {
			continue
		}
		if filter.releaseDate != nil && movie.ReleaseDate.Before(*filter.releaseDate) {
			continue
		}
		if filter.rating != nil && movie.Rating < *filter.rating {
			continue
		}
		res = append(res, movie)
	}

	return res
}
