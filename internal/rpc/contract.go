package rpc

import "context"

type DbRepo interface {
	SignupInsert(ctx context.Context, username string, email string, password string) error
}
