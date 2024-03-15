package queries

import (
	"context"
	"fmt"
	"time"
	"vk-backend/internal/domain"
)

const addMovieQuery = `INSERT INTO movies (title, description,release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id`

func (q *Queries) AddMovie(
	ctx context.Context,
	title string,
	description string,
	releaseDate time.Time,
	rating float64,
	actors []*domain.Actor,
) (*domain.Movie, error) {
	row := q.pool.QueryRow(ctx, addMovieQuery, title, description, releaseDate, rating)

	movie := &domain.Movie{
		Title:       title,
		Description: description,
		ReleaseDate: releaseDate,
		Rating:      rating,
		Actors:      actors,
	}
	if err := row.Scan(&movie.Id); err != nil {
		return nil, fmt.Errorf("failed to add movie: %w", err)
	}

	for _, actor := range actors {
		if err := q.AddActorToMovie(ctx, actor.Id, movie.Id); err != nil {
			return nil, fmt.Errorf("failed to add actor to movie: %w", err)
		}
	}

	return movie, nil
}

const getMovieByIdQuery = `SELECT id, title, description, release_date, rating FROM movies WHERE id = $1`
const getMovieActorsQuery = `SELECT actor_id FROM movie_actors WHERE movie_id = $1`

func (q *Queries) GetMovieById(ctx context.Context, id int) (*domain.Movie, error) {
	row := q.pool.QueryRow(ctx, getMovieByIdQuery, id)

	movie := &domain.Movie{}
	if err := row.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
		return nil, fmt.Errorf("failed to get movie by id: %w", err)
	}
	rows, err := q.pool.Query(ctx, getMovieActorsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie actors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var actorId int
		if err := rows.Scan(&actorId); err != nil {
			return nil, fmt.Errorf("failed to get movie actors: %w", err)
		}
		actor, err := q.GetActorById(ctx, actorId)
		if err != nil {
			return nil, fmt.Errorf("failed to get movie actors: %w", err)
		}
		movie.Actors = append(movie.Actors, actor)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to get movie actors: %w", rows.Err())
	}

	return movie, nil
}

const listMoviesQuery = `SELECT id, title, description, release_date, rating FROM movies`

func (q *Queries) ListMovies(ctx context.Context) ([]*domain.Movie, error) {
	rows, err := q.pool.Query(ctx, listMoviesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to list movies: %w", err)
	}
	defer rows.Close()

	var movies []*domain.Movie
	for rows.Next() {
		movie := &domain.Movie{}
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, fmt.Errorf("failed to list movies: %w", err)
		}
		movies = append(movies, movie)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to list movies: %w", rows.Err())
	}

	return movies, nil
}

const findMoviesByTitleQuery = `SELECT id, title, description, release_date, rating FROM movies WHERE title = $1`

func (q *Queries) FindMoviesByTitle(ctx context.Context, name string) ([]*domain.Movie, error) {
	rows, err := q.pool.Query(ctx, findMoviesByTitleQuery, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find movies by title: %w", err)
	}
	defer rows.Close()

	var movies []*domain.Movie
	for rows.Next() {
		movie := &domain.Movie{}
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, fmt.Errorf("failed to find movies by title: %w", err)
		}
		movies = append(movies, movie)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to find movies by title: %w", rows.Err())
	}

	return movies, nil
}

const findMoviesByActorNameQuery = `
SELECT m.id, m.title, m.description, m.release_date, m.rating
FROM movies AS m
JOIN movie_actors AS ma ON m.id = ma.movie_id
JOIN actors AS a ON ma.actor_id = a.id 
WHERE a.name = $1
`

func (q *Queries) FindMoviesByActorName(ctx context.Context, name string) ([]*domain.Movie, error) {
	rows, err := q.pool.Query(ctx, findMoviesByActorNameQuery, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find movies by actor name: %w", err)
	}
	defer rows.Close()

	var movies []*domain.Movie
	for rows.Next() {
		movie := &domain.Movie{}
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, fmt.Errorf("failed to find movies by actor name: %w", err)
		}
		movies = append(movies, movie)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to find movies by actor name: %w", rows.Err())
	}

	return movies, nil
}

const getMoviesWithActorQuery = `
SELECT m.id, m.title, m.description, m.release_date, m.rating
FROM movies AS m
JOIN movie_actors AS ma ON m.id = ma.movie_id
WHERE ma.actor_id = $1
`

func (q *Queries) GetMoviesWithActor(ctx context.Context, actorId int) ([]*domain.Movie, error) {
	rows, err := q.pool.Query(ctx, getMoviesWithActorQuery, actorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies with actor: %w", err)
	}
	defer rows.Close()

	var movies []*domain.Movie
	for rows.Next() {
		movie := &domain.Movie{}
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, fmt.Errorf("failed to get movies with actor: %w", err)
		}
		movies = append(movies, movie)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to get movies with actor: %w", rows.Err())
	}

	return movies, nil
}

const updateMovieQuery = `UPDATE movies SET title = $2, description = $3, release_date = $4, rating = $5 WHERE id = $1`

func (q *Queries) UpdateMovie(ctx context.Context, new *domain.Movie) error {
	if _, err := q.pool.Exec(ctx, updateMovieQuery, new.Id, new.Title, new.Description, new.ReleaseDate, new.Rating); err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	return nil
}

const deleteMovieQuery = `DELETE FROM movies WHERE id = $1`

func (q *Queries) DeleteMovie(ctx context.Context, id int) error {
	if _, err := q.pool.Exec(ctx, deleteMovieQuery, id); err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	return nil
}

const existsMovieQuery = `SELECT EXISTS(SELECT 1 FROM movies WHERE id = $1)`

func (q *Queries) MovieExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	if err := q.pool.QueryRow(ctx, existsMovieQuery, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check if movie exists: %w", err)
	}

	return exists, nil
}
