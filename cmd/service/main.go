package main

import (
	"github.com/Golang-Mentor-Education/auth/internal/rpc"
	"github.com/Golang-Mentor-Education/auth/pkg/auth"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	srv := rpc.NewService()

	s := grpc.NewServer()

	auth.RegisterAuthServiceServer(s, srv)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.Serve(lis)
}
