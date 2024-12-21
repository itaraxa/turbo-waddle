package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	e "github.com/itaraxa/turbo-waddle/internal/errors"
	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
	"github.com/itaraxa/turbo-waddle/internal/services"
)

type storager interface {
	services.UserStorager
	services.OrderStorager
	services.BalanceStorager
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
		startTime := time.Now()

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
			case errors.Is(err, e.ErrInvalidRequestFormat):
				http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
				l.Error("invalid request format", "error", err)
				return
			default:
				http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
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
		l.Info(`registration request completed`, `login`, user.Login, `password`, user.Password, `token`, token, "duration", time.Since(startTime))
	}
}

/*
Login - creates handler for authorisation existing user

Args:

	ctx context.Context
	l log.Logger
	s storager
	sk []byte: secret key for JWT-signing

Returns:

	http.HandlerFunc
*/
func Login(ctx context.Context, l log.Logger, s storager, sk []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info(`existing user authorisation request`)
		startTime := time.Now()

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
		l.Info("user data", "login", user.Login, "password", user.Password)
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

		// Call Authentication of existing user
		token, err := services.Authentication(ctx, l, s, user.Login, user.Password, sk)
		if err != nil {
			switch {
			case errors.Is(err, e.ErrInvalidLoginPassPair):
				http.Error(w, e.ErrInvalidLoginPassPair.Error(), e.ErrInvalidLoginPassPair.Code)
				l.Error("invalid pair of login and password", "login", user.Login, "password", user.Password)
				return
			case errors.Is(err, e.ErrUserNotFound):
				http.Error(w, e.ErrUserNotFound.Error(), e.ErrUserNotFound.Code)
				l.Error("user not found", "login", user.Login)
				return
			case errors.Is(err, e.ErrInvalidRequestFormat):
				http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
				l.Error("invalid request format", "error", err)
				return
			default:
				http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
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
		l.Info(`authorisation request completed`, `login`, user.Login, `password`, user.Password, `token`, token, "duration", time.Since(startTime))
	}
}

/*
PostOrders - create handler for adding new orders

Args:

	ctx context.Context
	l log.Logger
	s storager
	sk []byte: sekret key for signing/checking JWT token

Returns:

	http.HandlerFunc
*/
func PostOrders(ctx context.Context, l log.Logger, s storager, sk []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info("order adding request")
		startTime := time.Now()

		// Check authentication
		token := r.Header.Get("Autorisation")
		if token == "" {
			http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
			err := errors.Join(e.ErrUserIsNotauthenticated, errors.New("authetication token was not provided"))
			l.Error("user not authenticated", "error", err)
			return
		}
		login, err := services.CheckAuthentication(ctx, l, s, token, sk)
		if err != nil {
			http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
			err := errors.Join(e.ErrUserIsNotauthenticated, err)
			l.Error("user not authenticated", "error", err)
			return
		}
		l.Debug("request from user", "login", login)

		// Read data from request
		var buf bytes.Buffer
		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
			l.Error("cannot read from request body")
			return
		}

		// Check body
		orderNumber := buf.String()
		if len(orderNumber) == 0 {
			http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
			err := errors.Join(e.ErrInvalidRequestFormat, errors.New("empty order number"))
			l.Error("check order number error", "error", err)
			return
		}
		for ch := range orderNumber {
			if ch < '0' || ch > '9' {
				http.Error(w, e.ErrInvalidRequestFormat.Error(), e.ErrInvalidRequestFormat.Code)
				err := errors.Join(e.ErrInvalidRequestFormat, errors.New("non digit symbols in order number"))
				l.Error("check order number error", "error", err)
				return
			}
		}

		// Validate order number
		ok, err := services.ValidateOrderNumber(orderNumber, services.LUHN)
		if err != nil && errors.Is(err, e.ErrInvalidOrderNumberFormat) {
			http.Error(w, e.ErrInvalidOrderNumberFormat.Error(), e.ErrInvalidOrderNumberFormat.Code)
			l.Error("invalid order number format", "error", err)
			return
		}
		if err != nil {
			http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
			l.Error("validation order number error", "error", err)
			return
		}

		if !ok {
			http.Error(w, e.ErrInvalidOrderNumberFormat.Error(), e.ErrInvalidOrderNumberFormat.Code)
			l.Error("order number failed Luhn algorithm validation", "order number", orderNumber, "error", err)
			return
		}
		l.Debug("getted order number is valid", "order number", orderNumber)

		// Add order into storage
		err = services.LoadOrder(ctx, l, s, login, orderNumber)
		if err != nil {
			return
		}

		// Write and send response with Autorisation token
		w.Header().Set("Autorisation", token)
		w.WriteHeader(http.StatusAccepted)

		l.Info(`order adding request completed`, `duration`, time.Since(startTime))
	}
}

/*
GetOrders - return handler-function what returns list of orders

Args:

	ctx context.Context
	l log.Logger
	s storager
	sk []byte: sekret key for JWT-token

Returns:

	http.HandlerFunc
*/
func GetOrders(ctx context.Context, l log.Logger, s storager, sk []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info("getting orders request")
		startTIme := time.Now()

		// Check authentication
		token := r.Header.Get("Autorisation")
		if token == "" {
			http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
			err := errors.Join(e.ErrUserIsNotauthenticated, errors.New("authetication token was not provided"))
			l.Error("user not authenticated", "error", err)
			return
		}
		login, err := services.CheckAuthentication(ctx, l, s, token, sk)
		if err != nil {
			http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
			err := errors.Join(e.ErrUserIsNotauthenticated, err)
			l.Error("user not authenticated", "error", err)
			return
		}
		l.Debug("request from user", "login", login)

		// Getting orders from storage
		orders, err := services.GetOrders(ctx, l, s, login)
		if err != nil {
			switch {
			case errors.Is(err, e.ErrNoData):
				http.Error(w, e.ErrNoData.Error(), e.ErrNoData.Code)
				l.Error("no data for response", "error", err)
				return
			case errors.Is(err, e.ErrUserIsNotauthenticated):
				http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
				l.Error("user not authenticated", "error", err)
				return
			default:
				http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
				l.Error("internal server error", "error", err)
				return
			}
		}

		jsonData, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
			l.Error("internal server error: marshal data", "error", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
			l.Error("internal server error: writting body", "error", err)
			return
		}

		l.Info("getting orders reques completed", "duration", time.Since(startTIme))
	}
}

// Получение баланса
func GetBalance(ctx context.Context, l log.Logger, s storager, sk []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info("getting balance request")
		startTime := time.Now()

		// Check authentication
		token := r.Header.Get("Autorisation")
		if token == "" {
			http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
			err := errors.Join(e.ErrUserIsNotauthenticated, errors.New("authetication token was not provided"))
			l.Error("user not authenticated", "error", err)
			return
		}
		login, err := services.CheckAuthentication(ctx, l, s, token, sk)
		if err != nil {
			http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
			err := errors.Join(e.ErrUserIsNotauthenticated, err)
			l.Error("user not authenticated", "error", err)
			return
		}
		l.Debug("request from user", "login", login)

		// get balance
		balance, err := services.GetBalance(ctx, l, s, login)
		if err != nil {
			switch {
			case errors.Is(err, e.ErrUserIsNotauthenticated):
				http.Error(w, e.ErrUserIsNotauthenticated.Error(), e.ErrUserIsNotauthenticated.Code)
				err := errors.Join(e.ErrUserIsNotauthenticated, err)
				l.Error("user not authenticated", "error", err)
				return
			default:
				http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
				l.Error("internal server error", "error", err)
				return
			}
		}
		jsonData, err := json.Marshal(balance)
		if err != nil {
			http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
			l.Error("internal server error: marshal data", "error", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, e.ErrInternalServerError.Error(), e.ErrInternalServerError.Code)
			l.Error("internal server error: writting body", "error", err)
			return
		}

		l.Info("getting balance request completed", "duration", time.Since(startTime))
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
