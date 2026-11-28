package usecase

import (
	"context"

	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/domain"
)

// Login hold the repository
type Login struct {
	userRepository UserRepository
}

// NewLogin build the Login with repository
func NewLogin(userRepository UserRepository) *Login {
	return &Login{userRepository: userRepository}
}

// LoginInput is the raw data for registration.
type LoginInput struct {
	Email    string
	Password string
}

// GetByEmail validates the input,
// get the user by email and password checking
func (l *Login) GetByEmail(
	ctx context.Context,
	in LoginInput,
) (*domain.User, error) {
	email, err := domain.ParseLogin(in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	user, err := l.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.PasswordHash != in.Password {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}
