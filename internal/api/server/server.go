package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"vk-backend/internal/api/middleware"
	"vk-backend/internal/api/router"
	"vk-backend/internal/service/actor"
	"vk-backend/internal/service/movie"
)

type Server struct {
	srv *http.Server
}

func New(addr string, actorSrv *actor.Service, movieSrv *movie.Service, log *logrus.Logger) *Server {
	mux := middleware.LoggingMiddleware(router.New(actorSrv, movieSrv), log)
	srv := &http.Server{
		Addr:    ":" + addr,
		Handler: mux,
	}

	return &Server{srv: srv}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
