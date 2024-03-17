package handlers

import (
	"errors"
	"net/http"
	"vk-backend/internal/domain"
)

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
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Internal server error"))
	}
}
