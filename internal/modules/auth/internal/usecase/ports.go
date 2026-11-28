package usecase

import (
	"context"

	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/domain"
)

// UserRepository stores users port
type UserRepository interface {
	Create(ctx context.Context, u *domain.UserParams) (*domain.User, error)
}
