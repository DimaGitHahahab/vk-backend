package domain

import "errors"

var (
	ErrEmptyName       = errors.New("empty name")
	ErrFutureBirthDate = errors.New("birth date is in the future")
	ErrEmptyBirthDate  = errors.New("empty birth date")
	ErrInvalidGender   = errors.New("invalid gender")
	ErrActorNotExists  = errors.New("actor does not exist")
	ErrMovieNotExists  = errors.New("movie does not exist")

	ErrEmptyTitle         = errors.New("empty title")
	ErrTooLongTitle       = errors.New("title is too long")
	ErrEmptyDescription   = errors.New("empty description")
	ErrTooLongDescription = errors.New("description is too long")

	ErrInvalidRating       = errors.New("rating is invalid")
	ErrActorAlreadyInMovie = errors.New("actor is already in the movie")
	ErrEmptyReleaseDate    = errors.New("empty release date")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotExists     = errors.New("user does not exist")
	ErrInvalidLogin      = errors.New("invalid username or password")

	ErrEmptyPassword = errors.New("empty password")

	ErrNotAdmin = errors.New("not admin")
)
