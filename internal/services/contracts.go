package services

import (
	"context"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

type userStorager interface {
	AddNewUser(ctx context.Context, l log.Logger, login string, password string) (token string, err error)
	LoginUser(ctx context.Context, l log.Logger, login string, password string) (token string, err error)
}

type orderStorager interface{}

type balanceStorager interface{}
