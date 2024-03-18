package handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
	"vk-backend/internal/domain"
)

type AuthRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := &AuthRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid request body"))
		return
	}

	_, err := h.user.Register(r.Context(), req.Name, req.Password)
	if err != nil {
		h.HandleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

const oneMonth = time.Hour * 24 * 30

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := &AuthRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid request body"))
		return
	}

	u, err := h.user.GetUserByName(r.Context(), req.Name)
	if err != nil {
		h.HandleServiceError(w, err)
		return
	}

	if req.Password == "" {
		h.HandleServiceError(w, domain.ErrEmptyPassword)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		h.HandleServiceError(w, domain.ErrInvalidLogin)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  u.Id,
		"role": u.Role,
		"exp":  time.Now().Add(oneMonth).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		h.HandleServiceError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    t,
		Expires:  time.Now().Add(oneMonth),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)

}

func isAdminRole(r *http.Request) bool {
	role := r.Context().Value("user_role").(string)
	return role == "admin"
}
