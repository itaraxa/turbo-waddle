package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/itaraxa/turbo-waddle/internal/crypto"
	e "github.com/itaraxa/turbo-waddle/internal/errors"
	"github.com/itaraxa/turbo-waddle/internal/log"
)

// Зарегистрировать нового пользователя
func Registration(ctx context.Context, l log.Logger, us UserStorager, login string, password string, sk []byte) (token string, err error) {
	l.Info("registration new user", "login", login, "password", password)
	startTime := time.Now()

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
	l.Info("registration complited", "login", login, "token", token, "duration", time.Since(startTime))
	return
}

// Аутентификация пользователя
func Authentication(ctx context.Context, l log.Logger, us UserStorager, login string, password string, sk []byte) (token string, err error) {
	l.Info("authentication user", "login", login, "password", password)
	startTime := time.Now()

	salt, storedHash, err := us.GetUserHash(ctx, l, login)
	if err != nil {
		l.Error("authentication user error", "login", login, "error", err)
		return
	}
	l.Debug("hashes from storage", "salt", hex.EncodeToString(salt), "password_hash", hex.EncodeToString(storedHash))

	verifiedHash, err := crypto.GeneratePasswordWithSaltHash(salt, []byte(password))
	if err != nil {
		l.Error("generation hash for checking error", "login", login, "password", password, "error", err)
		err = errors.Join(err, e.ErrInternalServerError)
		return
	}

	if !bytes.Equal(storedHash, verifiedHash[:]) {
		l.Debug("hashes not equal", "wanted", hex.EncodeToString(storedHash), "got", hex.EncodeToString(verifiedHash[:]))
		l.Error("invalid password hash", "login", login, "password", password)
		err = e.ErrInvalidLoginPassPair
		return
	}

	token, err = crypto.CreateJWT(login, sk)
	if err != nil {
		l.Error("generating token for user error", "error", err)
		return "", errors.Join(e.ErrInternalServerError, err)
	}

	err = us.AddSession(ctx, l, login, token)
	if err != nil {
		l.Error("adding session into storage error", "login", login, "token", token, "error", err)
		return "", errors.Join(e.ErrInternalServerError, err)
	}
	l.Info("authentication complited", "login", login, "token", token, "duration", time.Since(startTime))
	return
}

// Проверка JWT и получение имени пользователя
func CheckAuthentication(ctx context.Context, l log.Logger, us UserStorager, token string, sk []byte) (login string, err error) {
	l.Info("authentication check", "token", token)

	valid, err := crypto.VerifyJWT(token, sk)
	if err != nil {
		l.Error("verify JWT error", "token", token, "error", err)
		return
	}
	if !valid {
		l.Error("invalid JWT token", "token", token)
		return "", e.ErrUserIsNotauthenticated
	}

	login, err = crypto.GetUsernameFromJWT(token, sk)
	if err != nil {
		l.Error("parse JWT error", "token", token, "error", err)
		return "", e.ErrUserIsNotauthenticated
	}

	l.Info("authentication check completed", "login", login)

	return
}
