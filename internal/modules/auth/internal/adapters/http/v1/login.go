package v1

import (
	"net/http"
	"time"

	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/usecase"
	"github.com/sohan-14/auth-system-go/internal/platform/httpx"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Login handles POST /auth/login. It only decodes and maps
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		h.ew.Write(w, r, err)
		return
	}

	user, err := h.login.GetByEmail(r.Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		h.ew.Write(w, r, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, loginResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}); err != nil {
		h.ew.Write(w, r, err)
	}
}
