package router

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"vk-backend/internal/api/handlers"
	"vk-backend/internal/api/middleware"
	"vk-backend/internal/service/actor"
	"vk-backend/internal/service/movie"
	"vk-backend/internal/service/user"
)

func New(actorSrv *actor.Service, movieSrv *movie.Service, user *user.Service, log *logrus.Logger) *http.ServeMux {
	h := handlers.New(*actorSrv, *movieSrv, *user)

	mux := http.NewServeMux()
	registerHandlerWithAuth(mux, "POST", "/actors", h.AddActorHandler, log)
	registerHandlerWithAuth(mux, "GET", "/actors", h.GetAllActorsHandler, log)
	registerHandlerWithAuth(mux, "GET", "/actors/{id}", h.GetActorHandler, log)
	registerHandlerWithAuth(mux, "PUT", "/actors/{id}", h.UpdateActorHandler, log)
	registerHandlerWithAuth(mux, "PATCH", "/actors/{id}", h.UpdateActorHandler, log)
	registerHandlerWithAuth(mux, "DELETE", "/actors/{id}", h.DeleteActorHandler, log)
	registerHandlerWithAuth(mux, "POST", "/movies", h.AddMovieHandler, log)
	registerHandlerWithAuth(mux, "POST", "/movies/{id}/actors", h.AddActorToMovieHandler, log)
	registerHandlerWithAuth(mux, "GET", "/movies", h.GetMoviesHandler, log)
	registerHandlerWithAuth(mux, "GET", "/movies/{id}", h.GetMovieHandler, log)
	registerHandlerWithAuth(mux, "PUT", "/movies/{id}", h.UpdateMovieHandler, log)
	registerHandlerWithAuth(mux, "PATCH", "/movies/{id}", h.UpdateMovieHandler, log)
	registerHandlerWithAuth(mux, "DELETE", "/movies/{id}", h.DeleteMovieHandler, log)

	// admin role must be given manually straight in db (task description), so there's no endpoint for that
	mux.Handle("/register", middleware.Logging(http.HandlerFunc(h.RegisterHandler), log))
	mux.Handle("/login", middleware.Logging(http.HandlerFunc(h.LoginHandler), log))

	return mux
}

func registerHandlerWithAuth(mux *http.ServeMux, method string, path string, h http.HandlerFunc, log *logrus.Logger) {
	mux.Handle(method+" "+path, middleware.Logging(middleware.RequireAuth(h), log))
}
