package accrual

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"

	"github.com/itaraxa/turbo-waddle/internal/log"
)

/*
order status from accrual system
*/
const (
	ORDER_REGISTERED = `REGISTERED`
	ORDER_PROCESSING = `PROCESSING`
	ORDER_INVALID    = `INVALID`
	ORDER_PROCESSED  = `PROCESSED`
)

var accrualStatuses = map[string]struct{}{
	ORDER_INVALID:    {},
	ORDER_PROCESSED:  {},
	ORDER_PROCESSING: {},
	ORDER_REGISTERED: {},
}

type ClientAccrual struct {
	httpClient      *resty.Client
	accrualEndpoint string
}

func NewAccrualSystem(accrualEndpoint string) *ClientAccrual {
	client := resty.New()
	return &ClientAccrual{
		httpClient:      client,
		accrualEndpoint: accrualEndpoint,
	}
}

type OrderAccrual struct {
	OrderNumber  string          `json:"order"`
	OrderStatus  string          `json:"status"`
	OrderAccrual decimal.Decimal `json:"accrual"`
}

/*
GetOrderAccrual executes GET-request to accrual system. Response will be checked for returned status

Args:

	ctx context.Context
	l log.Logger
	orderNumber string

Returns:

	status string
	accrual decimal.Decimal
	err error
*/
func (ca *ClientAccrual) GetOrderAccrual(ctx context.Context, l log.Logger, orderNumber string) (status string, accrual decimal.Decimal, err error) {
	url := strings.Join([]string{ca.accrualEndpoint, `api`, `orders`, orderNumber}, `/`)
	l.Debug("requst to accrual system", "url", url)
	var oa OrderAccrual
	resp, err := ca.httpClient.R().SetResult(&oa).Get(url)
	if err != nil {
		l.Error("cannot do request to accrual system", "error", err)
		err = errors.Join(ErrInternalServerError, err)
		return
	}

	l.Debug("response from accrual system", "code", resp.StatusCode(), "body", string(resp.Body()))
	switch resp.StatusCode() {
	case 200:
		l.Debug("data from responce", "order", oa.OrderNumber, "status", oa.OrderStatus, "accrual", oa.OrderAccrual)
		if _, ok := accrualStatuses[oa.OrderStatus]; !ok {
			err = ErrUnknownStatus
			return
		}
		status = oa.OrderStatus
		accrual = oa.OrderAccrual
		return
	case 204:
		err = ErrOrderDoesNotRegistered
		l.Error("response error", "error", err)
		return
	case 429:
		err1 := ErrServerRequestLimitExceeded
		retryAfter, err2 := strconv.Atoi(resp.Header().Get("Retry-After"))
		if err2 != nil {
			l.Error("cannot get retry interval from header")
			err = errors.Join(err, err2)
			return
		}
		err1.RetryAfter = retryAfter
		l.Error("response error", "error", err1, "retry-after", err1.RetryAfter)
		err = err1
		return
	default:
		err = ErrInternalServerError
		l.Error("response error", "error", err)
		return
	}
}
