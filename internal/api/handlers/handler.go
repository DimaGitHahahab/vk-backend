package handlers

import (
	"net/http"
	"vk-backend/internal/service/actor"
	"vk-backend/internal/service/movie"
)

type Handler struct {
	act actor.Service
	mov movie.Service
}

func New(act actor.Service, mov movie.Service) *Handler {
	return &Handler{
		act: act,
		mov: mov,
	}
}

func (h *Handler) PingHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("pong"))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
