package app

import (
	sl "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itaraxa/turbo-waddle/internal/config"
	"github.com/itaraxa/turbo-waddle/internal/log"
	"github.com/itaraxa/turbo-waddle/internal/storage"
	"github.com/itaraxa/turbo-waddle/internal/tranposrt/rest"
)

type logger interface {
	Error(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

type storager interface{}

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

func NewServerApp(config *config.GopherMartConfig) *ServerApp {
	l, err := log.NewZapLogger(config.LogLevel)
	if err != nil {
		sl.Fatal(err)
	}

	storage, err := storage.NewStorage(config.DSN)
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
	sa.log.Info(`server started`,
		`version`, version,
	)

	// setup router
	sa.r.Use(rest.Logger(),
		rest.ChekcUser(),
		rest.CheckRequest(),
		rest.Decompress(),
		rest.Compress(),
	)
	sa.r.Post(`/api/user/register`, rest.Register())
	sa.r.Post(`/api/user/login`, rest.Login())
	sa.r.Post(`/api/user/orders`, rest.PostOrders())
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

	defer sa.log.Info(`server stopped`)
}
