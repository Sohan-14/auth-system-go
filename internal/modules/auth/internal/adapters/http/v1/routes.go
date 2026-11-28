package v1

import "github.com/sohan-14/auth-system-go/internal/platform/router"

// Routes registers the v1 endpoints; r supplies the version prefix.
func (h *Handler) Routes(r *router.Router) {
	r.HandleFunc("POST /register", h.Register)
}
