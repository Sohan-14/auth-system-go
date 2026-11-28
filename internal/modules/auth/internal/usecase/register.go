package usecase

import (
	"context"

	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/domain"
)

// Register hold the
type Register struct {
	userRepository UserRepository
}

// NewRegister build the register with port
func NewRegister(userRepository UserRepository) *Register {
	return &Register{userRepository: userRepository}
}

// RegisterInput is the raw data for registration.
type RegisterInput struct {
	Email    string
	Password string
}

// Create validates the input, hashes the password and stores the user.
func (r *Register) Create(
	ctx context.Context,
	in RegisterInput,
) (*domain.User, error) {
	email, err := domain.ParseRegistration(in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepository.Create(ctx, &domain.UserParams{
		Email:        email,
		PasswordHash: in.Password,
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
