package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sohan-14/auth-system-go/internal/apperr"
	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/domain"
)

// GetByEmail find the user by email.
func (repository *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	user, err := repository.q.GetUserByEmail(ctx, email)

	switch {
	case err == nil:
		return &domain.User{
			ID:           user.ID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt.Time,
			UpdatedAt:    user.UpdatedAt.Time,
		}, nil
	case errors.Is(err, pgx.ErrNoRows):
		return nil, domain.ErrUserNotFound
	default:
		return nil, apperr.Internal("get user by email", err)
	}
}
