package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Pinkman-77/records-restapi/pkg/models"
	"github.com/Pinkman-77/records-restapi/pkg/storage"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveUser inserts a new user into the users table
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt := `INSERT INTO users(email, password_hash) VALUES($1, $2) RETURNING id`

	var id int64
	err := s.db.QueryRowContext(ctx, stmt, email, passHash).Scan(&id)
	if err != nil {
		// Handle unique constraint error
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // PostgreSQL unique violation
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User retrieves a user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	stmt := `SELECT id, email, password_hash FROM users WHERE email = $1`

	var user models.User
	err := s.db.QueryRowContext(ctx, stmt, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	stmt := `SELECT is_admin FROM users WHERE id = $1`	

	var isAdmin bool

	err := s.db.QueryRowContext(ctx, stmt, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)

	}

	return isAdmin, nil
}

// App retrieves an application by ID
func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"

	stmt := `SELECT id, name, secret FROM apps WHERE id = $1`

	var app models.App
	err := s.db.QueryRowContext(ctx, stmt, id).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}


