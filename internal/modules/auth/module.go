package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/sohan-14/auth-system-go/internal/modules/auth/internal/adapters/http/v1"
	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/adapters/postgres"
	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/usecase"
	"github.com/sohan-14/auth-system-go/internal/platform/httpx"
	"github.com/sohan-14/auth-system-go/internal/platform/router"
)

// Module is the public auth entry point.
type Module struct {
	v1Handler *v1.Handler
}

// New wires adapters, use cases and the HTTP handler.
func New(
	pool *pgxpool.Pool,
	errorWriter *httpx.ErrorWriter,
) *Module {
	userRepository := postgres.NewUserRepository(pool)
	registerUsecase := usecase.NewRegister(userRepository)

	return &Module{
		v1Handler: v1.NewHandler(
			registerUsecase,
			errorWriter,
		),
	}
}

// RegisterRoutes registers the v1 routes
func (m *Module) RegisterRoutes(router *router.Router) {
	authRouter := router.Group("/api/v1/auth")
	m.v1Handler.Routes(authRouter)
}
