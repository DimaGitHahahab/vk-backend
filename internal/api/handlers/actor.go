package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"vk-backend/internal/domain"
)

type ActorRequest struct {
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type ActorDTO struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

func (h *Handler) AddActorHandler(writer http.ResponseWriter, request *http.Request) {
	act := &ActorRequest{}
	if err := json.NewDecoder(request.Body).Decode(act); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid request body"))
		return
	}

	if !isAdminRole(request) {
		h.HandleServiceError(writer, domain.ErrNotAdmin)
		return
	}

	actor, err := h.act.AddActor(request.Context(), act.Name, genderStringToInt(act.Gender), act.BirthDate)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	dto := actorToDTO(actor)

	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(dto); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Internal server error"))
		return
	}
}

func (h *Handler) UpdateActorHandler(writer http.ResponseWriter, request *http.Request) {
	act := &ActorRequest{}
	if err := json.NewDecoder(request.Body).Decode(act); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid request body"))
		return
	}

	if !isAdminRole(request) {
		h.HandleServiceError(writer, domain.ErrNotAdmin)
		return
	}

	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid actor id"))
		return
	}

	actor, err := h.act.GetActorById(request.Context(), id)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	if act.Name != "" {
		actor.Name = act.Name
	}
	if act.Gender != "" {
		actor.Gender = genderStringToInt(act.Gender)
	}
	if !act.BirthDate.IsZero() {
		actor.BirthDate = act.BirthDate
	}

	if err := h.act.UpdateActor(request.Context(), actor); err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetActorHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid actor id"))
		return
	}

	actor, err := h.act.GetActorById(request.Context(), id)
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	dto := actorToDTO(actor)

	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(dto); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Internal server error"))
		return
	}
}

func (h *Handler) GetAllActorsHandler(writer http.ResponseWriter, request *http.Request) {
	actors, err := h.act.ListActors(request.Context())
	if err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	var dtos []ActorDTO
	for _, a := range actors {
		dtos = append(dtos, actorToDTO(a))
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

func (h *Handler) DeleteActorHandler(writer http.ResponseWriter, request *http.Request) {

	if !isAdminRole(request) {
		h.HandleServiceError(writer, domain.ErrNotAdmin)
		return
	}

	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("Invalid actor id"))
		return
	}

	if err := h.act.DeleteActor(request.Context(), id); err != nil {
		h.HandleServiceError(writer, err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)

}

func actorToDTO(actor *domain.Actor) ActorDTO {
	a := ActorDTO{
		Id:        actor.Id,
		Name:      actor.Name,
		BirthDate: actor.BirthDate,
	}
	switch actor.Gender {
	case 0:
		a.Gender = "unknown"
	case 1:
		a.Gender = "male"
	case 2:
		a.Gender = "female"
	default:
		a.Gender = "not applicable"
	}

	return a
}
func genderStringToInt(g string) int {
	switch g {
	case "unknown":
		return 0
	case "male":
		return 1
	case "female":
		return 2
	case "not applicable":
		return 9
	}

	return -1
}
