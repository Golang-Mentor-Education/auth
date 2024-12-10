package rpc

import "context"

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type DbRepo interface {
	SignupInsert(ctx context.Context, username string, email string, password string) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
