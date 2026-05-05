package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MiltonJ23/Daedalus/AuthService/internal/domain"
	"github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	stmt, prepareContextErr := r.db.PrepareContext(ctx, query)
	if prepareContextErr != nil {
		return prepareContextErr
	}
	defer stmt.Close()

	_, execContextErr := stmt.ExecContext(ctx, user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	if execContextErr != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(execContextErr, &mysqlErr) && mysqlErr.Number == 1062 {
			return errors.New("user already exists")
		}
		return execContextErr
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = ?`

	var user domain.User

	fetchingUserErr := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if fetchingUserErr != nil {
		if errors.Is(fetchingUserErr, sql.ErrNoRows) {
			return nil, nil // It simply means the user doesn't exist , it is not an error
		}
		return nil, fetchingUserErr
	}
	return &user, nil
}
