package v1

import (
	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/usecase"
	"github.com/sohan-14/auth-system-go/internal/platform/httpx"
)

// Handler serves the v1 auth endpoints.
type Handler struct {
	register *usecase.Register
	login    *usecase.Login
	ew       *httpx.ErrorWriter
}

// NewHandler build the handler for v1 with usecases and error writer
func NewHandler(
	register *usecase.Register,
	login *usecase.Login,
	ew *httpx.ErrorWriter,
) *Handler {
	return &Handler{register: register, login: login, ew: ew}
}
