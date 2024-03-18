package handlers

import (
	"errors"
	"net/http"
	"vk-backend/internal/domain"
	"vk-backend/internal/service/actor"
	"vk-backend/internal/service/movie"
	"vk-backend/internal/service/user"
)

type Handler struct {
	act  actor.Service
	mov  movie.Service
	user user.Service
}

func New(act actor.Service, mov movie.Service, user user.Service) *Handler {
	return &Handler{
		act:  act,
		mov:  mov,
		user: user,
	}
}

func (h *Handler) HandleServiceError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyName):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Name cannot be empty"))
	case errors.Is(err, domain.ErrFutureBirthDate):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Birth date cannot be in the future"))
	case errors.Is(err, domain.ErrEmptyBirthDate):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Birth date cannot be empty"))
	case errors.Is(err, domain.ErrInvalidGender):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid gender. Can be 'unknown', 'male', 'female', 'not applicable'"))
	case errors.Is(err, domain.ErrActorNotExists):
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("Actor does not exist"))
	case errors.Is(err, domain.ErrMovieNotExists):
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("Movie does not exist"))
	case errors.Is(err, domain.ErrEmptyTitle):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Title cannot be empty"))
	case errors.Is(err, domain.ErrTooLongTitle):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Title is too long"))
	case errors.Is(err, domain.ErrEmptyDescription):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Description cannot be empty"))
	case errors.Is(err, domain.ErrTooLongDescription):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Description is too long"))
	case errors.Is(err, domain.ErrInvalidRating):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Rating is invalid"))
	case errors.Is(err, domain.ErrActorAlreadyInMovie):
		writer.WriteHeader(http.StatusConflict)
		_, _ = writer.Write([]byte("Actor is already in the movie"))
	case errors.Is(err, domain.ErrEmptyReleaseDate):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Release date cannot be empty"))
	case errors.Is(err, domain.ErrUserAlreadyExists):
		writer.WriteHeader(http.StatusConflict)
		_, _ = writer.Write([]byte("User already exists"))
	case errors.Is(err, domain.ErrUserNotExists):
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("User does not exist"))
	case errors.Is(err, domain.ErrInvalidLogin):
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte("Invalid username or password"))
	case errors.Is(err, domain.ErrEmptyPassword):
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Password cannot be empty"))
	case errors.Is(err, domain.ErrNotAdmin):
		writer.WriteHeader(http.StatusForbidden)
		_, _ = writer.Write([]byte("not allowed"))
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Internal server error"))
	}
}
