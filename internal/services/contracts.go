package services

import (
	"context"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

type UserStorager interface {
	AddNewUser(ctx context.Context, l log.Logger, login string, hash []byte, salt []byte, token string) (err error)
	LoginUser(ctx context.Context, l log.Logger, login string, password string, token string) (err error)
	GetUserHash(ctx context.Context, l log.Logger, login string) (salt, hash []byte, err error)
	AddSession(ctx context.Context, l log.Logger, login string, token string) (err error)
}

type StoragerUtilites interface {
	HelthCheck(ctx context.Context, l log.Logger) (err error)
}

type orderStorager interface{}

type balanceStorager interface{}
