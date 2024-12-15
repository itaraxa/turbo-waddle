package services

import (
	"context"
	"errors"

	"github.com/itaraxa/turbo-waddle/internal/crypto"
	e "github.com/itaraxa/turbo-waddle/internal/errors"
	"github.com/itaraxa/turbo-waddle/internal/log"
)

// Зарегистрировать нового пользователя
func Registration(ctx context.Context, l log.Logger, us UserStorager, login string, password string, sk []byte) (token string, err error) {
	l.Info("registration new user", "login", login, "password", password)

	salt, err := crypto.GenerateSalt(32)
	if err != nil {
		return "", errors.Join(ErrUserRegistration, err)
	}
	hash, err := crypto.GeneratePasswordWithSaltHash(salt, []byte(password))
	if err != nil {
		return "", errors.Join(ErrUserRegistration, err)
	}
	token, err = crypto.CreateJWT(login, sk)
	if err != nil {
		l.Error("generating token for new user error", "error", err)
		return "", errors.Join(e.ErrInternalServerError, err)
	}

	err = us.AddNewUser(ctx, l, login, hash[:], salt, token)
	if err != nil {
		l.Error("adding new user error", "login", login, "error", err)
		return "", errors.Join(ErrUserRegistration, err)
	}
	l.Info("registration complited", "login", login, "token", token)
	return
}

// Аутентификация пользователя
func Authentication(ctx context.Context, l log.Logger, us UserStorager, login string, password string) (token string, err error) {
	l.Info("authentication user", "login", login, "password", password)

	token, err = crypto.GenerateToken64()
	if err != nil {
		l.Error("generating toke for new user error", "error", err)
		return "", errors.Join(e.ErrInternalServerError, err)
	}

	err = us.LoginUser(ctx, l, login, password, token)
	if err != nil {
		l.Error("authentication user error", "login", login, "error", err)
		return "", errors.Join(ErrUserAuthentication, err)
	}
	l.Info("authentication complited", "login", login, "token", token)
	return
}
