package rpc

import (
	"context"
	"github.com/Golang-Mentor-Education/auth/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Service struct {
	auth.UnimplementedAuthServiceServer
	dbRepo DbRepo
}

func NewService(dbR DbRepo) *Service {
	return &Service{
		dbRepo: dbR,
	}
}

func (s *Service) Login(ctx context.Context, in *auth.LoginIn) (*auth.LoginOut, error) {
	return &auth.LoginOut{
		Token: in.Username + in.Password,
	}, nil
}

func (s *Service) Signup(ctx context.Context, in *auth.SignupIn) (*auth.SignupOut, error) {
	if in.Email == "" || in.Password == "" || in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password or username is empty")
	}
	err := s.dbRepo.SignupInsert(ctx, in.Username, in.Email, in.Password)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.SignupOut{
		Success: true,
	}, nil
}
