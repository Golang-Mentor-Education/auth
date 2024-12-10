package repository

import (
	"context"
	"fmt"
	"github.com/Golang-Mentor-Education/auth/internal/rpc"
	sq "github.com/Masterminds/squirrel"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	postgresPort     = "3011"
	postgresUser     = "master"
	postgresPassword = "master"
	postgresDb       = "master"
	postgresHost     = "localhost"
)

type Repository struct {
	conn *sqlx.DB
}

func NewRepository() *Repository {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", postgresUser, postgresPassword, postgresDb, postgresHost, postgresPort)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return &Repository{conn: conn}
}

func (r *Repository) SignupInsert(ctx context.Context, username string, email string, password string) error {
	query, args, err := sq.Insert("participant").
		Columns("username", "email", "password").
		Values(username, email, password).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*rpc.User, error) {
	query, args, err := sq.Select("id", "username", "email", "password").
		From("participant").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var u rpc.User
	err = r.conn.QueryRowxContext(ctx, query, args...).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*rpc.User, error) {
	query, args, err := sq.Select("id", "username", "email", "password").
		From("participant").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var u rpc.User
	err = r.conn.QueryRowxContext(ctx, query, args...).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &u, nil
}
