package accrual

import (
	"context"
	"net/http"

	"github.com/itaraxa/turbo-waddle/internal/config"

	"github.com/shopspring/decimal"
)

type logger interface{}

type accrualer interface {
	GetOrder() string
	GetStatus() string
	GetAccrual() decimal.Decimal
}

type ClientAccrual struct {
	logger     logger
	config     *config.GopherMartConfig
	httpClient *http.Client
}

func NewClientAccrual(l logger, c *config.GopherMartConfig, h *http.Client) *ClientAccrual {
	return &ClientAccrual{
		logger:     l,
		config:     c,
		httpClient: h,
	}
}

func (ca *ClientAccrual) GetOrderAccrual(ctx context.Context, l logger, orderNumber string) (accrualer, error) {
	return nil, nil
}
