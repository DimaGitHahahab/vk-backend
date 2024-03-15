package domain

import "errors"

var (
	ErrEmptyName       = errors.New("empty name")
	ErrFutureBirthDate = errors.New("birth date is in the future")
	ErrEmptyBirthDate  = errors.New("empty birth date")
	ErrActorNotExists  = errors.New("actor does not exist")
	ErrMovieNotExists  = errors.New("movie does not exist")

	ErrEmptyTitle         = errors.New("empty title")
	ErrTooLongTitle       = errors.New("title is too long")
	ErrEmptyDescription   = errors.New("empty description")
	ErrTooLongDescription = errors.New("description is too long")

	ErrorRatingInvalid     = errors.New("rating is invalid")
	ErrActorAlreadyInMovie = errors.New("actor is already in the movie")
)
