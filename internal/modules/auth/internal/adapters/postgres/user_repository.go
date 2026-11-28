package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sohan-14/auth-system-go/internal/modules/auth/internal/adapters/postgres/db"
)

// UserRepository stores users in PostgreSQL.
type UserRepository struct {
	q *db.Queries
}

// NewUserRepository accepts a pool and build the UserRepository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: db.New(pool)}
}
