package postgres

import (
	"context"

	"github.com/sohan-14/auth-system-go/internal/apperr"
	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/adapters/postgres/db"
	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/domain"
)

// Create inserts user, mapping a duplicate email to domain.ErrEmailTaken.
func (repository *UserRepository) Create(
	ctx context.Context,
	input *domain.UserParams,
) (*domain.User, error) {
	user, err := repository.q.CreateUser(ctx, db.CreateUserParams{
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
	})

	switch {
	case err == nil:
		return &domain.User{
			ID:    user.ID,
			Email: user.Email,
		}, nil
	case isUniqueViolation(err):
		return nil, domain.ErrEmailTaken
	default:
		return nil, apperr.Internal("create user", err)
	}

}
