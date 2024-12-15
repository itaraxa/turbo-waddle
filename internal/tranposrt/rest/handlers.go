package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	e "github.com/itaraxa/turbo-waddle/internal/errors"
	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
	"github.com/itaraxa/turbo-waddle/internal/services"
)

type storager interface {
	services.UserStorager
}

/*
Register - creates handler for registration new user

Args:

	ctx context.Context
	l log.Logger
	s storager
	sk []byte: secret key for JWT-signing

Returns:

	http.HandlerFunc
*/
func Register(ctx context.Context, l log.Logger, s storager, sk []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info(`new user registration request`)

		// Read data from request
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
			l.Error("cannot read from request body")
			return
		}
		var user models.User
		if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
			http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
			l.Error("cannot unmarshal request body")
			return
		}

		// Check data from request
		l.Info("New user data", "login", user.Login, "password", user.Password)
		if user.Login == "" {
			http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
			l.Error("getted empty user login")
			return
		}
		if user.Password == "" {
			http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
			l.Error("getted empty user password")
			return
		}

		// Call Registration of new user
		token, err := services.Registration(ctx, l, s, user.Login, user.Password, sk)
		if err != nil {
			switch {
			case errors.Is(err, e.ErrLoginIsAlreadyUsed):
				http.Error(w, e.ErrLoginIsAlreadyUsed.Error(), e.ErrLoginIsAlreadyUsed.Code)
				l.Error("login already used", "login", user.Login, "error", err)
				return
			case errors.Is(err, e.ErrInternalServerError):
				http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
				l.Error("internal server error", "error", err)
				return
			case errors.Is(err, e.ErrInvalidRequestFormat):
				http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
				l.Error("invalid request format", "error", err)
				return
			default:
				http.Error(w, e.ErrLoginIsAlreadyUsed.Error(), e.ErrLoginIsAlreadyUsed.Code)
				l.Error("internal server error", "error", err)
				return
			}
		}

		// Write and send response with Autorisation token
		w.Header().Set("Autorisation", token)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(""))
		if err != nil {
			http.Error(w, "cannot write to response body", http.StatusNoContent)
			l.Error("cannot write to response body", "error", err)
			return
		}
	}
}

// Авторизация пользователя
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Добавление заказа
func PostOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Получение списка заказов
func GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Получение баланса
func GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Запрос на вывод бонусов
func WithdrawRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Получение списка выводов
func GetWithdrawls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
