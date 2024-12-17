package app

import (
	"context"
	sl "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/itaraxa/turbo-waddle/internal/config"
	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/services"
	"github.com/itaraxa/turbo-waddle/internal/storage"
	"github.com/itaraxa/turbo-waddle/internal/tranposrt/rest"
	"github.com/itaraxa/turbo-waddle/internal/version"
)

type logger interface {
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

type storager interface {
	services.UserStorager
	services.HealthCheckStorager
	services.OrderStorager
}

type router interface {
	http.Handler
	Use(middlewares ...func(next http.Handler) http.Handler)
	Get(pattern string, handlerFn http.HandlerFunc)
	Post(pattern string, handlerFn http.HandlerFunc)
	// ListenAndServe(addr string, handler http.Handler) error
}

type ServerApp struct {
	log     logger
	storage storager
	r       router
	config  *config.GopherMartConfig
}

func NewServerApp(ctx context.Context, config *config.GopherMartConfig) *ServerApp {
	l, err := log.NewZapLogger(config.LogLevel)
	if err != nil {
		sl.Fatal(err)
	}
	l.Info(`server started`,
		`app version`, version.ServerApp,
		`database schema version`, version.Database,
	)
	l.Debug(`server configuration`,
		`endpoint`, config.Endpoint,
		`database`, config.DSN,
		`accrual system address`, config.AccrualSystemAddress,
		`log level`, config.LogLevel,
		`show version`, config.ShowVersion,
	)

	storage, err := storage.NewStorage(ctx, l, config.DSN)
	if err != nil {
		sl.Fatal(err)
	}

	return &ServerApp{
		log:     l,
		storage: storage,
		r:       chi.NewRouter(),
		config:  config,
	}
}

func (sa *ServerApp) Run() {
	sa.log.Debug(`Test DEBUG message`)
	sa.log.Info(`Test INFO message`)
	sa.log.Warn(`Test WARN message`)
	sa.log.Error(`Test ERROR message`)

	stopAppChannel := make(chan os.Signal, 1)
	signal.Notify(stopAppChannel, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(cancel context.CancelFunc) {
		defer cancel()
		<-stopAppChannel
		sa.log.Info(`stopping gophermat app`, `reason`, `getted interrupt signal from OS`)
		close(stopAppChannel)
	}(cancel)

	// database healthcheck
	go func(ctx context.Context, l logger, stopCh chan os.Signal, s storager, period time.Duration) {
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := s.HealthCheck(ctx, l); err != nil {
					l.Error("database healthcheck failed", "error", err)
				} else {
					l.Info("database connection ok")
				}
			case <-stopCh:
				return
			}
		}
	}(ctx, sa.log, stopAppChannel, sa.storage, 10*time.Second)

	// setup router
	sa.r.Use(rest.Logger(),
		rest.ChekcUser(),
		rest.CheckRequest(),
		rest.Decompress(),
		rest.Compress(),
	)
	sa.r.Post(`/api/user/register`, rest.Register(ctx, sa.log, sa.storage, sa.config.SecretKey))
	sa.r.Post(`/api/user/login`, rest.Login(ctx, sa.log, sa.storage, sa.config.SecretKey))
	sa.r.Post(`/api/user/orders`, rest.PostOrders(ctx, sa.log, sa.storage, sa.config.SecretKey))
	sa.r.Get(`/api/user/orders`, rest.GetOrders())
	sa.r.Get(`/api/user/balance`, rest.GetBalance())
	sa.r.Post(`/api/user/balance/withdraw`, rest.WithdrawRequest())
	sa.r.Get(`/api/user/withdrawals`, rest.GetWithdrawls())

	server := &http.Server{
		Addr:    sa.config.Endpoint,
		Handler: sa.r,
	}

	// start router
	go func() {
		sa.log.Info("start router")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			sa.log.Fatal("router error", "err", err.Error())
		}
	}()

	// stop router
	<-ctx.Done()
	sa.log.Info(`stopping router`, `reason`, `context was cancele`)
	ctxWithTimeout, cancelWithTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelWithTimeout()
	err := server.Shutdown(ctxWithTimeout)
	if err != nil {
		sa.log.Fatal(`stopping router`, `error`, err)
	}
	sa.log.Info(`router stopped`)

	defer sa.log.Info(`server stopped`)
}
