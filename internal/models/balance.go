package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}

type Withdraw struct {
	Order       string
	Sum         decimal.Decimal
	ProcessedAt time.Time
}

func (b *Balance) String() string {
	return fmt.Sprintf("Current ballnce: %s, Withdrawn: %s", b.Current.String(), b.Withdrawn.String())
}
