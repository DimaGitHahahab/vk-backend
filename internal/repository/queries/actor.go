package queries

import (
	"context"
	"fmt"
	"time"
	"vk-backend/internal/domain"
)

const insertActorQuery = `INSERT INTO actors (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id`

func (q *Queries) AddActor(ctx context.Context, name string, gender int, birthDate time.Time) (*domain.Actor, error) {
	row := q.pool.QueryRow(ctx, insertActorQuery, name, gender, birthDate)

	actor := &domain.Actor{
		Name:      name,
		Gender:    gender,
		BirthDate: birthDate,
	}
	if err := row.Scan(&actor.Id); err != nil {
		return nil, fmt.Errorf("failed to insert actor: %w", err)
	}

	return actor, nil
}

const insertActorToMovieQuery = `INSERT INTO movie_actors (actor_id, movie_id) VALUES ($1, $2)`

func (q *Queries) AddActorToMovie(ctx context.Context, actorId int, movieId int) error {
	if _, err := q.pool.Exec(ctx, insertActorToMovieQuery, actorId, movieId); err != nil {
		return fmt.Errorf("failed to insert actor to movie: %w", err)
	}

	return nil
}

const selectActorQuery = `SELECT name, gender, birth_date FROM actors WHERE id = $1`

func (q *Queries) GetActorById(ctx context.Context, id int) (*domain.Actor, error) {
	row := q.pool.QueryRow(ctx, selectActorQuery, id)

	actor := &domain.Actor{Id: id}
	if err := row.Scan(&actor.Name, &actor.Gender, &actor.BirthDate); err != nil {
		return nil, fmt.Errorf("failed to get actor: %w", err)
	}

	return actor, nil
}

const selectActorsByMovieIdQuery = `
SELECT actors.id, actors.name, actors.gender, actors.birth_date
FROM actors
JOIN movie_actors ON actors.id = movie_actors.actor_id
WHERE movie_actors.movie_id = $1
`

func (q *Queries) GetActorsByMovieId(ctx context.Context, movieId int) ([]*domain.Actor, error) {
	rows, err := q.pool.Query(ctx, selectActorsByMovieIdQuery, movieId)
	if err != nil {
		return nil, fmt.Errorf("failed to select actors by movie id: %w", err)
	}

	defer rows.Close()

	var actors []*domain.Actor
	for rows.Next() {
		actor := &domain.Actor{}
		if err := rows.Scan(ctx, selectActorsByMovieIdQuery, &actor.Id, &actor.Name, &actor.Gender, &actor.BirthDate); err != nil {
			return nil, fmt.Errorf("failed to get actors by movie id: %w", err)
		}
		actors = append(actors, actor)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to get actors by movie id: %w", rows.Err())
	}

	return actors, nil
}

const selectAllActorsQuery = `SELECT id, name, gender, birth_date FROM actors`

func (q *Queries) ListActors(ctx context.Context) ([]*domain.Actor, error) {
	rows, err := q.pool.Query(ctx, selectActorQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to select all actors: %w", err)
	}

	defer rows.Close()

	var actors []*domain.Actor
	for rows.Next() {
		actor := &domain.Actor{}
		if err := rows.Scan(ctx, selectAllActorsQuery, &actor.Id, &actor.Name, &actor.Gender, &actor.Gender); err != nil {
			return nil, fmt.Errorf("failed to list all the actors: %w", err)
		}
		actors = append(actors, actor)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to list all the actors: %w", rows.Err())
	}

	return actors, nil
}

const updateActorQuery = `UPDATE actors SET name = $2, gender = $3, birth_date = $4 WHERE id = $1`

func (q *Queries) UpdateActor(ctx context.Context, new *domain.Actor) error {
	if _, err := q.pool.Exec(ctx, updateActorQuery, new.Id, new.Name, new.Gender, new.BirthDate); err != nil {
		return fmt.Errorf("failed to update actor: %w", err)
	}

	return nil
}

const deleteActorQuery = `DELETE FROM actors WHERE id = $1`

func (q *Queries) DeleteActor(ctx context.Context, id int) error {
	if _, err := q.pool.Exec(ctx, deleteActorQuery, id); err != nil {
		return fmt.Errorf("failed to delete actor: %w", err)
	}

	return nil
}

const existsActorQuery = `SELECT EXISTS(SELECT 1 FROM actors WHERE id = $1)`

func (q *Queries) ActorExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	if err := q.pool.QueryRow(ctx, existsActorQuery, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check if actor exists: %w", err)
	}

	return exists, nil
}
