package movie

import (
	"strings"
	"time"
	"vk-backend/internal/domain"
)

type Filter struct {
	name        *string
	releaseDate *time.Time
	rating      *float64
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) WithTitle(title string) *Filter {
	f.name = &title
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
		if filter.name != nil &&
			!strings.Contains(strings.ToLower(movie.Title), strings.ToLower(*filter.name)) &&
			!searchActor(movie.Actors, *filter.name) {
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

func searchActor(actors []*domain.Actor, name string) bool {
	for _, actor := range actors {
		if strings.Contains(strings.ToLower(actor.Name), strings.ToLower(name)) {
			return true
		}
	}
	return false
}
