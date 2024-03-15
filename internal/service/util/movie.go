package util

import (
	"sort"
	"strings"
	"time"
	"vk-backend/internal/domain"
)

type Filter struct {
	Title       *string
	ReleaseDate *time.Time
	Rating      *float64
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) WithTitle(title string) *Filter {
	f.Title = &title
	return f
}

func (f *Filter) WithReleaseDate(releaseDate time.Time) *Filter {
	f.ReleaseDate = &releaseDate
	return f
}

func (f *Filter) WithRating(rating float64) *Filter {
	f.Rating = &rating
	return f
}

type SortBy int

const (
	SortByRating SortBy = iota
	SortByReleaseDate
	SortByTitle

	defaultSort = SortByRating
)

func FilterMovies(movies []*domain.Movie, filter *Filter) []*domain.Movie {
	res := make([]*domain.Movie, 0, len(movies))
	for _, movie := range movies {
		if filter.Title != nil && !strings.Contains(movie.Title, *filter.Title) {
			continue
		}
		if filter.ReleaseDate != nil && movie.ReleaseDate != *filter.ReleaseDate {
			continue
		}
		if filter.Rating != nil && movie.Rating < *filter.Rating {
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
		return domain.ErrorRatingInvalid
	}

	return nil
}
