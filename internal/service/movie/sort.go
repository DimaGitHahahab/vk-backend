package movie

import (
	"sort"
	"vk-backend/internal/domain"
)

type SortBy int

const (
	SortByRating SortBy = iota
	SortByReleaseDate
	SortByTitle

	DefaultSort = SortByRating
)

func SortMovies(movies []*domain.Movie, sorting SortBy) []*domain.Movie {
	sortFunc := func(i, j int) bool {
		switch sorting {
		case SortByRating:
			return movies[i].Rating > movies[j].Rating
		case SortByReleaseDate:
			return movies[j].ReleaseDate.Before(movies[i].ReleaseDate)
		case SortByTitle:
			return movies[i].Title < movies[j].Title
		default:
			return movies[i].Rating < movies[j].Rating
		}
	}
	sort.Slice(movies, sortFunc)

	return movies
}
