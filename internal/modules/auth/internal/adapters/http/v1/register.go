package v1

import (
	"net/http"
	"time"

	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/usecase"
	"github.com/sohan-14/auth-system-go/internal/platform/httpx"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Register handles POST /auth/register. It only decodes and maps
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		h.ew.Write(w, r, err)
		return
	}

	user, err := h.register.Create(r.Context(), usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		h.ew.Write(w, r, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusCreated, registerResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}); err != nil {
		h.ew.Write(w, r, err)
	}
}
