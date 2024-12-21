package services

import (
	"context"
	"sync"
	"time"

	"github.com/itaraxa/turbo-waddle/internal/client/accrual"
	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/models"
	"github.com/shopspring/decimal"
)

const (
	JobCount              = 10
	WorkerCount           = 3
	AccrualRequestTimeout = 3 * time.Second
	UpdateOrderDBTimeout  = 3 * time.Second
)

// Добавление запроса на расчет бонусов за заказ
func LoadOrder(ctx context.Context, l log.Logger, os OrderStorager, login string, order string) (err error) {
	err = os.LoadOrder(ctx, l, login, order)
	if err != nil {
		l.Error("loading order into storage error", "error", err)
	}
	return
}

// Получение информации о бонусах пользователя
func GetOrders(ctx context.Context, l log.Logger, os OrderStorager, login string) (orders []models.Order, err error) {
	orders, err = os.GetOrders(ctx, l, login)
	if err != nil {
		l.Error("getting orders from storage error")
	}
	return
}

// Обновление статусов из Accrual System
// используется workerPool
type Job struct {
	ID    int64
	Order string
}

type Result struct {
	JobID   int64
	Order   string
	Status  string
	Accrual decimal.Decimal
	Err     error
}

func worker(ctx context.Context, l log.Logger, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, as *accrual.ClientAccrual) {
	defer wg.Done()

LOOP:
	for {
		select {
		case job := <-jobs:
			ctxWithTimeout, cancel := context.WithTimeout(ctx, AccrualRequestTimeout)
			status, accrual, err := as.GetOrderAccrual(ctxWithTimeout, l, job.Order)
			// TO-DO: add retry
			if err != nil {
				l.Error("getting order accrual error", "job_id", job.ID, "order", job.Order, "worker_id", id, "error", err)
				cancel()
				continue
			}
			l.Debug("getting order accrual completed", "job_id", job.ID, "order", job.Order, "worker_id", id, "status", status, "accrual", accrual)

			results <- Result{JobID: job.ID, Order: job.Order, Status: status, Accrual: accrual, Err: err}
			cancel()
		case <-ctx.Done():
			break LOOP
		}
	}
}

func collectResults(ctx context.Context, l log.Logger, results <-chan Result, wg *sync.WaitGroup, os OrderStorager) {
	defer wg.Done()
	for result := range results {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, UpdateOrderDBTimeout)

		err := os.UpdateOrder(ctxWithTimeout, l, result.Order, result.Status, result.Accrual)
		if err != nil {
			l.Error("updating order in storage error", "order", result.Order, "status", result.Status, "accrual", result.Accrual, "error", err)
			cancel()
			continue
		}
		l.Debug("updating order in storage completed", "order", result.Order, "status", result.Status, "accrual", result.Accrual)
		cancel()
	}
}

/*
AccrualUpdate checks order statuses in storage and requests for updates to Accrual System

Args:

	ctx context.Context
	l log.Logger
	os OrderStorager
	accrualEndpoint string
*/
func AccrualUpdate(ctx context.Context, l log.Logger, os OrderStorager, accrualEndpoint string) {
	l.Info("accrual updating started", "accrual system", accrualEndpoint)
	startTime := time.Now()

	as := accrual.NewAccrualSystem(accrualEndpoint)
	jobs := make(chan Job, JobCount)
	results := make(chan Result, JobCount)
	var wg sync.WaitGroup

	wg.Add(WorkerCount)
	for w := 1; w <= WorkerCount; w++ {
		go worker(ctx, l, w, jobs, results, &wg, as)
		l.Debug("worker created", "worker id", w)
	}

	wg.Add(1)
	go collectResults(ctx, l, results, &wg, os)

	wg.Add(1)
	ticker1 := time.NewTicker(1000 * time.Millisecond)
	defer ticker1.Stop()
	go func(ctx context.Context, ticker *time.Ticker, l log.Logger, jobs chan<- Job, wg *sync.WaitGroup) {
		defer wg.Done()
		var jobId uint64 = 0
	LOOP:
		for {
			select {
			case <-ticker.C:
				notProcessedOrders, err := os.GetNotProcessedOrders(ctx, l)
				if err != nil {
					l.Error("getting orders from storage error", "error", err)
					continue
				}
				for _, order := range notProcessedOrders {
					job := Job{
						ID:    int64(jobId),
						Order: order.Order,
					}
					jobId++
					jobs <- job
				}
			case <-ctx.Done():
				break LOOP
			}
		}
		close(jobs)
	}(ctx, ticker1, l, jobs, &wg)

	ticker5 := time.NewTicker(5 * time.Second)
	defer ticker5.Stop()
	for {
		select {
		case <-ticker5.C:
			l.Debug("accrual update working")
		case <-ctx.Done():
			l.Info("accrual updating stopping", "reason", "context canceled")
			wg.Wait()
			l.Info("accrual updating stopped", "duration", time.Since(startTime))
			return
		}
	}
}
