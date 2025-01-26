package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/Golang-Mentor-Education/auth/internal/config"
	"github.com/Golang-Mentor-Education/auth/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Тест на успешный логин по username
func TestService_Login_SuccessByUsername(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDbRepo(ctrl)

	// Фиктивная конфигурация
	cfg := &config.Config{
		Platform: config.Platform{
			Token: "supersecretkey",
		},
	}

	// Создаём сервис
	srv := NewService(cfg, mockDbRepo)

	// Входные данные
	in := &auth.LoginIn{
		Username: "alex",
		Password: "123456",
		// Email: ""
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// Мокаем возвращаемого пользователя
	user := &User{
		ID:       1,
		Username: "alex",
		Email:    "alex@example.com",
		Password: string(hash),
	}

	// При вызове GetUserByUsername → вернём user, nil
	mockDbRepo.EXPECT().GetUserByUsername(gomock.Any(), "alex").Return(user, nil).Times(1)

	// Вызываем Login
	out, err := srv.Login(context.Background(), in)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Token, "должен вернуться не пустой токен")
}

// Тест на отсутствие пароля
func TestService_Login_NoPassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDbRepo(ctrl)
	cfg := &config.Config{
		Platform: config.Platform{
			Token: "secret",
		},
	}

	srv := NewService(cfg, mockDbRepo)

	// Пароль пуст
	in := &auth.LoginIn{
		Username: "alex",
	}

	out, err := srv.Login(context.Background(), in)

	assert.Error(t, err, "ожидаем ошибку, т.к. пароль пуст")
	assert.Nil(t, out, "ответ должен быть nil")
}

// Тест на успешный логин по email
func TestService_Login_SuccessByEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDbRepo(ctrl)
	cfg := &config.Config{
		Platform: config.Platform{
			Token: "supersecretkey",
		},
	}
	srv := NewService(cfg, mockDbRepo)

	in := &auth.LoginIn{
		Email:    "alex@example.com",
		Password: "xyz987",
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &User{
		ID:       2,
		Username: "alex-user",
		Email:    "alex@example.com",
		Password: string(hash),
	}

	mockDbRepo.EXPECT().GetUserByEmail(gomock.Any(), "alex@example.com").Return(user, nil).Times(1)

	out, err := srv.Login(context.Background(), in)
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Token)
}

// Тест на неверный пароль
func TestService_Login_InvalidCredentials(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDbRepo(ctrl)
	cfg := &config.Config{
		Platform: config.Platform{
			Token: "secret",
		},
	}
	srv := NewService(cfg, mockDbRepo)

	in := &auth.LoginIn{
		Username: "someone",
		Password: "wrong",
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("rightpass"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &User{
		ID:       100,
		Username: "someone",
		Email:    "some@where.com",
		Password: string(hash),
	}

	mockDbRepo.EXPECT().GetUserByUsername(gomock.Any(), "someone").Return(user, nil).Times(1)

	out, err := srv.Login(context.Background(), in)
	assert.Error(t, err, "Should fail on bad password")
	assert.Nil(t, out)
}
func TestService_Signup_Success(t *testing.T) {
	t.Parallel()

	// Готовим окружение
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDbRepo(ctrl)
	cfg := &config.Config{
		Platform: config.Platform{
			Token: "someKey", // Нам может не пригодиться, но пусть будет
		},
	}
	srv := NewService(cfg, mockDbRepo)

	// Входные данные
	in := &auth.SignupIn{
		Username: "alex",
		Email:    "alex@example.com",
		Password: "123456",
	}

	// Настраиваем мок: при вызове SignupInsert → nil (ошибок нет)
	mockDbRepo.EXPECT().
		SignupInsert(gomock.Any(), "alex", "alex@example.com", gomock.Any()).
		Return(nil).
		Times(1)

	// Вызываем Signup
	out, err := srv.Signup(context.Background(), in)

	// Проверяем результат
	assert.NoError(t, err, "Signup должен быть без ошибок")
	assert.True(t, out.Success, "Поле Success=true")
}

func TestService_Signup_DBError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDbRepo(ctrl)
	cfg := &config.Config{Platform: config.Platform{Token: "key"}}
	srv := NewService(cfg, mockDbRepo)

	in := &auth.SignupIn{
		Username: "alex",
		Email:    "alex@example.com",
		Password: "123456",
	}

	// Допустим, база вернет ошибку
	mockDbRepo.EXPECT().
		SignupInsert(gomock.Any(), "alex", "alex@example.com", gomock.Any()).
		Return(fmt.Errorf("db insertion failed")).
		Times(1)

	out, err := srv.Signup(context.Background(), in)
	assert.Error(t, err)
	assert.Nil(t, out)
}
