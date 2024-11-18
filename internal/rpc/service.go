package rpc

import (
	"context"
	"github.com/Golang-Mentor-Education/auth/pkg/auth"
)

type Service struct {
	auth.UnimplementedAuthServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Login(ctx context.Context, in *auth.LoginIn) (*auth.LoginOut, error) {
	return &auth.LoginOut{
		Token: in.Username + in.Password,
	}, nil
}
