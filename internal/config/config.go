package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/itaraxa/turbo-waddle/internal/version"
)

const (
	LOG_DEBUG = `DEBUG`
	LOG_INFO  = `INFO`
	LOG_WARN  = `WARN`
	LOG_ERROR = `ERROR`
)

var (
	ErrConfig                    = errors.New("Config: gophermart app configuration error")
	ErrParseFlag                 = errors.New("parseFlags: error command lina flags parsing")
	ErrParseEnv                  = errors.New("parseEnv: error environment variables parsing")
	ErrEmptyEndpoint             = errors.New("validateConfig: empty endpoint address")
	ErrEmptyDSN                  = errors.New("validateConfig: empty DSN address")
	ErrEmptyAccrualSystemAddress = errors.New("validateConfig: empty Accrual system address")
	ErrUnknownLogLevel           = errors.New("validateConfig: unknown log level")
)

type GopherMartConfig struct {
	Endpoint             string
	DSN                  string // postgres://gophermart:\!qaz2wsx@localhost:5432/gophermart
	AccrualSystemAddress string
	LogLevel             string
	ShowVersion          bool
}

func NewGopherMartConfig() *GopherMartConfig {
	return &GopherMartConfig{
		LogLevel: LOG_DEBUG,
	}
}

/*
Config runs a sequence of reading and validating configuration

returns:

	err error
*/
func (gmc *GopherMartConfig) Config() (err error) {
	err = gmc.parseFlags()
	if err != nil {
		return errors.Join(err, ErrConfig)
	}
	err = gmc.parseEnv()
	if err != nil {
		return errors.Join(err, ErrConfig)
	}
	err = gmc.validateConfig()
	if err != nil {
		return errors.Join(err, ErrConfig)
	}
	return
}

/*
parseFlags parses application flags and update configuration fields accordingly

Returns:

	err error
*/
func (gmc *GopherMartConfig) parseFlags() (err error) {
	flag.BoolVar(&gmc.ShowVersion, `v`, false, `Show version and exit`)
	flag.StringVar(&gmc.LogLevel, `l`, `INFO`, `Set log level: INFO, DEBUG, etc.`)

	flag.StringVar(&gmc.Endpoint, `a`, ``, `HTTP-server endpoint. Environment variable RUN_ADDRESS`)
	flag.StringVar(&gmc.DSN, `d`, ``, `Database source name. Environment variable DATABASE_URI`)
	flag.StringVar(&gmc.AccrualSystemAddress, `r`, ``, `Accrual system addres. Environment variable ACCRUAL_SYSTEM_ADDRESS`)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Gophermart app version: %s\n\rDatabase schema version: %d\n\rUsage of %s\n\r",
			version.ServerApp,
			version.Database,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	err = flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		err = errors.Join(ErrParseFlag, err)
		return
	}
	return
}

/*
parseEnv parses environment variables and update configuration fields accordingly

Returns:

	err error
*/
func (gmc *GopherMartConfig) parseEnv() (err error) {
	addressServer, ok := os.LookupEnv(`RUN_ADDRESS`)
	if ok {
		gmc.Endpoint = addressServer
	}
	dbDSN, ok := os.LookupEnv(`DATABASE_URI`)
	if ok {
		gmc.DSN = dbDSN
	}
	accrualSystemAddress, ok := os.LookupEnv(`ACCRUAL_SYSTEM_ADDRESS`)
	if ok {
		gmc.AccrualSystemAddress = accrualSystemAddress
	}

	return
}

/*
validateConfig checks specified configuration parametres.
Return error, if configuration parametr is empty, or contains
unknown value

Returns:

	err error
*/
func (gmc *GopherMartConfig) validateConfig() (err error) {
	if len(gmc.Endpoint) == 0 {
		err = errors.Join(err, ErrEmptyEndpoint)
	}

	if len(gmc.DSN) == 0 {
		err = errors.Join(err, ErrEmptyDSN)
	}

	if len(gmc.AccrualSystemAddress) == 0 {
		err = errors.Join(err, ErrEmptyAccrualSystemAddress)
	}

	logLevels := map[string]struct{}{
		LOG_DEBUG: {},
		LOG_INFO:  {},
		LOG_WARN:  {},
		LOG_ERROR: {},
	}
	if _, ok := logLevels[gmc.LogLevel]; !ok {
		err = errors.Join(err, ErrUnknownLogLevel)
	}

	return
}
