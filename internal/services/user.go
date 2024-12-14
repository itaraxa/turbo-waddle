package services

import (
	"context"
	"errors"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

// Зарегистрировать нового пользователя
func Registration(ctx context.Context, l log.Logger, us userStorager, login string, password string) (token string, err error) {
	l.Info("registration new user", "login", login, "password", password)
	token, err = us.AddNewUser(ctx, l, login, password)
	if err != nil {
		l.Error("registration user error", "login", login, "error", err)
		return "", errors.Join(ErrUserRegistration, err)
	}
	l.Info("registration complited", "login", login, "token", token)
	return
}

// Аутентификация пользователя
func Authentication(ctx context.Context, l log.Logger, us userStorager, login string, password string) (token string, err error) {
	l.Info("authentication user", "login", login, "password", password)
	token, err = us.LoginUser(ctx, l, login, password)
	if err != nil {
		l.Error("authentication user error", "login", login, "error", err)
		return "", errors.Join(ErrUserAuthentication, err)
	}
	l.Info("authentication complited", "login", login, "token", token)
	return
}
