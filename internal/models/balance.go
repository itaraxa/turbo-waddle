package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Ballance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}

type Withdraw struct {
	Order       string
	Sum         decimal.Decimal
	ProcessedAt time.Time
}
