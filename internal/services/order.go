package services

import (
	"context"

	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
)

// Добавление запроса на расчет бонусов за заказ
func LoadOrder(ctx context.Context, l log.Logger, os OrderStorager, login string, order string) (err error) {
	err = os.LoadOrder(ctx, l, login, order)
	return nil
}

// Получение информации о бонусах пользователя
func GetOrders(ctx context.Context, l log.Logger, os OrderStorager, login string) (orders []models.Order, err error) {

	return
}
