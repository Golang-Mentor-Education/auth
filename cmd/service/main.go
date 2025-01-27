package main

import (
	"fmt"
	"github.com/Golang-Mentor-Education/auth/internal/config"
	"log"
	"net"

	"github.com/Golang-Mentor-Education/auth/internal/repository"
	"github.com/Golang-Mentor-Education/auth/internal/rpc"
	"github.com/Golang-Mentor-Education/auth/pkg/auth"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.NewConfig()

	repo := repository.NewRepository(cfg)

	srv := rpc.NewService(cfg, repo)

	s := grpc.NewServer()

	auth.RegisterAuthServiceServer(s, srv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Auth running on :%s", cfg.Service.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
