package router

import (
	"net/http"
	"vk-backend/internal/api/handlers"
	"vk-backend/internal/service/actor"
	"vk-backend/internal/service/movie"
)

func New(actorSrv *actor.Service, movieSrv *movie.Service) *http.ServeMux {

	h := handlers.New(*actorSrv, *movieSrv)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /actors", h.AddActorHandler)
	mux.HandleFunc("GET /actors", h.GetAllActorsHandler)
	mux.HandleFunc("GET /actors/{id}", h.GetActorHandler)
	mux.HandleFunc("PUT /actors/{id}", h.UpdateActorHandler)
	mux.HandleFunc("PATCH /actors/{id}", h.UpdateActorHandler)
	mux.HandleFunc("DELETE /actors/{id}", h.DeleteActorHandler)
	mux.HandleFunc("POST /movies", h.AddMovieHandler)
	mux.HandleFunc("POST /movies/{id}/actors", h.AddActorToMovieHandler)
	mux.HandleFunc("GET /movies", h.GetAllMoviesHandler)
	mux.HandleFunc("GET /movies/{id}", h.GetMovieHandler)
	mux.HandleFunc("PUT /movies/{id}", h.UpdateMovieHandler)
	mux.HandleFunc("PATCH /movies/{id}", h.UpdateMovieHandler)
	mux.HandleFunc("DELETE /movies/{id}", h.DeleteMovieHandler)

	return mux
}
