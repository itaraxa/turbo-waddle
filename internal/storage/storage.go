package storage

import (
	"context"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

type Storager interface {
	UserStorager
}

type UserStorager interface {
	AddNewUser(ctx context.Context, l log.Logger, login string, password string) (token string, err error)
	LoginUser(ctx context.Context, l log.Logger, login string, password string) (token string, err error)
}

type Storage struct {
	dsn string
}

func NewStorage(dsn string) (*Storage, error) {
	return &Storage{
		dsn: dsn,
	}, nil
}
