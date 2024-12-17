package rpc

import (
	"context"
	"github.com/Golang-Mentor-Education/auth/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestServer_Login(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockDbRepo := NewMockDbRepo(ctrl)

	t.Run("success_username", func(t *testing.T) {

		in := &auth.LoginIn{
			Username: "alex",
			Password: "123456",
			Email:    "alex@gmail.com",
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		user := &User{
			Password: string(hash),
			ID:       1,
		}

		mockDbRepo.EXPECT().GetUserByUsername(gomock.Any(), in.Username).Return(user, nil)

		srv := NewService(mockDbRepo)
		token, err := srv.generateJWT(user.ID)
		assert.NoError(t, err)

		out, err := srv.Login(ctx, in)

		assert.NoError(t, err)
		assert.Equal(t, token, out.Token)
	})

	t.Run("no_password", func(t *testing.T) {
		in := &auth.LoginIn{
			Username: "alex",
		}

		srv := NewService(mockDbRepo)

		_, err := srv.Login(ctx, in)

		assert.Error(t, err)
	})
}
