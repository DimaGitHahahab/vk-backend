package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/internal/service/movie"
)

type MovieRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
	Actors      []int     `json:"actors"` // actor ids
}

type MovieDTO struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ReleaseDate time.Time  `json:"release_date"`
	Rating      float64    `json:"rating"`
	Actors      []ActorDTO `json:"actors"`
}

func (h *Handler) AddMovieHandler(writer http.ResponseWriter, request *http.Request) {
	mov := &MovieRequest{}
	if err := json.NewDecoder(request.Body).Decode(mov); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid request body"))
		return
	}

	actors := make([]*domain.Actor, 0, len(mov.Actors))
	for _, id := range mov.Actors {
		actor, err := h.act.GetActorById(request.Context(), id)
		if err != nil {
			h.HandleServiceError(writer, err)
			return
		}
		actors = append(actors, actor)
	}

	movie, err := h.mov.AddMovie(request.Context(), mov.Title, mov.Description, mov.ReleaseDate, mov.Rating, actors)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	dto := movieToDTO(movie)

	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(dto); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Internal server error"))
		return
	}
}

type AddActorToMovieRequest struct {
	ActorId int `json:"actor_id"`
}

func (h *Handler) AddActorToMovieHandler(writer http.ResponseWriter, request *http.Request) {
	req := &AddActorToMovieRequest{}
	if err := json.NewDecoder(request.Body).Decode(req); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid request body"))
		return
	}

	movieId, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid movie id"))
		return
	}

	err = h.mov.AddActorToMovie(request.Context(), req.ActorId, movieId)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
func (h *Handler) GetMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid movie id"))
		return
	}

	movie, err := h.mov.GetMovieById(request.Context(), id)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	dto := movieToDTO(movie)

	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(dto); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Internal server error"))
		return
	}
}

func (h *Handler) GetAllMoviesHandler(writer http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	sortParam := queryParams.Get("sort")
	var sort movie.SortBy
	switch sortParam {
	case "rating":
		sort = movie.SortByRating
	case "release_date":
		sort = movie.SortByReleaseDate
	case "title":
		sort = movie.SortByTitle
	default:
		sort = movie.DefaultSort
	}

	filter := movie.NewFilter()
	if title := queryParams.Get("title"); title != "" {
		filter = filter.WithTitle(title)
	}
	if releaseDate := queryParams.Get("release_date"); releaseDate != "" {
		parsedDate, err := time.Parse(time.DateOnly, releaseDate)
		if err == nil {
			filter = filter.WithReleaseDate(parsedDate)
		}
	}
	if rating := queryParams.Get("rating"); rating != "" {
		parsedRating, err := strconv.ParseFloat(rating, 64)
		if err == nil {
			filter = filter.WithRating(parsedRating)
		}
	}

	movies, err := h.mov.ListMovies(request.Context(), filter, sort)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	var dtos []MovieDTO
	for _, m := range movies {
		dtos = append(dtos, movieToDTO(m))
	}

	if len(dtos) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(dtos); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Internal server error"))
		return
	}
}

func (h *Handler) UpdateMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid movie id"))
		return
	}

	mov := &MovieRequest{}
	if err := json.NewDecoder(request.Body).Decode(mov); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid request body"))
		return
	}

	actors := make([]*domain.Actor, 0, len(mov.Actors))
	for _, id := range mov.Actors {
		actor, err := h.act.GetActorById(request.Context(), id)
		if err != nil {
			h.HandleServiceError(writer, err)
			return
		}
		actors = append(actors, actor)
	}

	oldMovie, err := h.mov.GetMovieById(request.Context(), id)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	if mov.Title != "" {
		oldMovie.Title = mov.Title
	}
	if mov.Description != "" {
		oldMovie.Description = mov.Description
	}

	if !mov.ReleaseDate.IsZero() {
		oldMovie.ReleaseDate = mov.ReleaseDate
	}

	if mov.Rating != 0 {
		oldMovie.Rating = mov.Rating
	}

	if len(actors) > 0 {
		oldMovie.Actors = actors
	}

	err = h.mov.UpdateMovie(request.Context(), oldMovie)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid movie id"))
		return
	}

	err = h.mov.DeleteMovie(request.Context(), id)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)

}

func movieToDTO(m *domain.Movie) MovieDTO {
	actors := make([]ActorDTO, 0, len(m.Actors))
	for _, a := range m.Actors {
		actors = append(actors, actorToDTO(a))
	}

	return MovieDTO{
		Id:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		ReleaseDate: m.ReleaseDate,
		Rating:      m.Rating,
		Actors:      actors,
	}
}
