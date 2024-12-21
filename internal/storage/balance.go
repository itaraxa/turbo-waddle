package storage

import (
	"context"
	"time"

	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
)

func (s *Storage) GetBalance(ctx context.Context, l log.Logger, login string) (bal models.Balance, err error) {
	l.Debug("getting balance from storage")
	startTime := time.Now()

	bal, err = s.PostgresRepository.GetBalance(ctx, l, login)
	if err != nil {
		l.Error("getting balance from storage error", "error", err)
		return
	}

	l.Debug("getting balance from storage completed", "duration", time.Since(startTime), "balance", bal.String(), "login", login)
	return
}
