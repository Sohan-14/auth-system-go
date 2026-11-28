package usecase

import (
	"context"

	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/domain"
)

// UserRepository stores users port
type UserRepository interface {
	Create(ctx context.Context, userParams *domain.UserParams) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
