package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	Number     string
	Status     string
	Accrual    decimal.Decimal
	UploadedAt time.Time
}
