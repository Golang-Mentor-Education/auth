package rpc

import (
	"context"
	"log"
	"time"

	"github.com/Golang-Mentor-Education/auth/internal/config"
	"github.com/Golang-Mentor-Education/auth/pkg/auth"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	auth.UnimplementedAuthServiceServer
	dbRepo    DbRepo
	jwtSecret []byte
}

func NewService(cfg *config.Config, dbR DbRepo) *Service {
	return &Service{
		dbRepo:    dbR,
		jwtSecret: []byte(cfg.Platform.Token),
	}
}

func (s *Service) Login(ctx context.Context, in *auth.LoginIn) (*auth.LoginOut, error) {
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	var user *User
	var err error

	// Проверяем что указал клиент
	if in.Username != "" {
		user, err = s.dbRepo.GetUserByUsername(ctx, in.Username)
	} else if in.Email != "" {
		user, err = s.dbRepo.GetUserByEmail(ctx, in.Email)
	} else {
		return nil, status.Error(codes.InvalidArgument, "username or email is required")
	}

	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Сравниваем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Генерируем JWT
	tokenString, err := s.generateJWT(user.ID)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &auth.LoginOut{
		Token: tokenString,
	}, nil
}

// JWT:
func (s *Service) generateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // токен на 24 часа
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) Signup(ctx context.Context, in *auth.SignupIn) (*auth.SignupOut, error) {
	if in.Email == "" || in.Password == "" || in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password or username is empty")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Insert hash into database instead password
	err = s.dbRepo.SignupInsert(ctx, in.Username, in.Email, string(hash))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.SignupOut{
		Success: true,
	}, nil
}
