package services

import (
	"context"

	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"

	"github.com/shopspring/decimal"
)

// Получить количество имеющихся бонусных баллов и уже потраченных
func GetBalance(ctx context.Context, l log.Logger, bs BalanceStorager, login string) (bal models.Balance, err error) {
	bal, err = bs.GetBalance(ctx, l, login)
	if err != nil {
		l.Error("getting balance from storage error", "error", err)
	}

	return
}

// Потратить балы на указанный заказ
func PerformWithdraw(ctx context.Context, l log.Logger, bs BalanceStorager, user string, order string, sum decimal.Decimal) error {

	return nil
}

// Получить информацию о совершенных тратах баллов
func GetWithdrawals(ctx context.Context, l log.Logger, bs BalanceStorager, user string) (w []models.Withdraw, err error) {
	return
}
