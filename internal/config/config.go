package config

type GopherMartConfig struct {
	Endpoint             string
	DatabaseDSN          string // postgres://gophermart:\!qaz2wsx@localhost:5433/gophermart
	AccrualSystemAddress string
	LogLevel             string
	ShowVersion          bool
}

func NewGopherMartConfig() *GopherMartConfig {
	return &GopherMartConfig{
		Endpoint:    `localhost:8080`,
		LogLevel:    `INFO`,
		ShowVersion: false,
	}
}
