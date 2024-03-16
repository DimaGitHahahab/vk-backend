package util

import (
	"sort"
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

type SortBy int

const (
	SortByRating SortBy = iota
	SortByReleaseDate
	SortByTitle

	DefaultSort = SortByRating
)

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

func SortMovies(movies []*domain.Movie, sorting SortBy) []*domain.Movie {
	sortFunc := func(i, j int) bool {
		switch sorting {
		case SortByRating:
			return movies[i].Rating < movies[j].Rating
		case SortByReleaseDate:
			return movies[i].ReleaseDate.Before(movies[j].ReleaseDate)
		case SortByTitle:
			return movies[i].Title < movies[j].Title
		default:
			return movies[i].Rating < movies[j].Rating
		}
	}
	sort.Slice(movies, sortFunc)

	return movies
}

func ValidateMovieData(title, description string, rating float64) error {
	if title == "" {
		return domain.ErrEmptyTitle
	}
	if len(title) > 100 {
		return domain.ErrTooLongTitle
	}
	if description == "" {
		return domain.ErrEmptyDescription
	}
	if len(description) > 1000 {
		return domain.ErrTooLongDescription
	}
	if rating < 0 || rating > 10 {
		return domain.ErrInvalidRating
	}

	return nil
}
