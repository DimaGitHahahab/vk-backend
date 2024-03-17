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
	mux.HandleFunc("GET /ping", h.PingHandler)
	mux.HandleFunc("POST /actors", h.AddActorHandler)
	mux.HandleFunc("GET /actors", h.GetAllActorsHandler)
	mux.HandleFunc("GET /actors/{id}", h.GetActorHandler)
	mux.HandleFunc("PUT /actors/{id}", h.UpdateActorHandler)
	mux.HandleFunc("PATCH /actors/{id}", h.UpdateActorHandler)
	mux.HandleFunc("DELETE /actors/{id}", h.DeleteActorHandler)

	return mux
}
