package storage

import (
	"context"
	"errors"
	"time"

	"github.com/itaraxa/turbo-waddle/internal/database/postgres"
	"github.com/itaraxa/turbo-waddle/internal/log"
)

type Storager interface {
	UserStorager
	HelthCheck(ctx context.Context, l log.Logger) (err error)
}

type UserStorager interface {
	AddNewUser(ctx context.Context, l log.Logger, login string, hash []byte, salt []byte, token string) (err error)
	GetUserHash(ctx context.Context, l log.Logger, login string) (salt, hash []byte, err error)
	LoginUser(ctx context.Context, l log.Logger, login string, password string, token string) (err error)
	AddSession(ctx context.Context, l log.Logger, login string, token string) (err error)
}

type Storage struct {
	*postgres.PostgresRepository
	dsn string
}

func NewStorage(ctx context.Context, l log.Logger, dsn string) (*Storage, error) {
	l.Info("Start creating new storage", "database source name", dsn)
	pr, err := postgres.NewPostgresRepository(ctx, l, dsn)
	if err != err {
		l.Error("Creating new storage error", "error", err)
		return nil, errors.Join(ErrCreateNewStorage, err)
	}
	l.Info("Strorrage created", "database source name", dsn)

	return &Storage{
		pr,
		dsn,
	}, nil
}

func (s *Storage) HelthCheck(ctx context.Context, l log.Logger) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err = s.PingContext(ctxWithTimeout); err != nil {
		return err
	}
	return
}
