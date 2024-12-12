package services

import (
	"context"

	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
)

// Добавление запроса на расчет бонусов за заказ
func LoadOrder(ctx context.Context, l log.Logger, os orderStorager, user string, order string) error {

	return nil
}

// Получение информации о бонусах пользователя
func GetOrders(ctx context.Context, l log.Logger, os orderStorager, user string) (orders []models.Order, err error) {

	return
}
