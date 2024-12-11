package main

import (
	"github.com/itaraxa/turbo-waddle/internal/app"
	"github.com/itaraxa/turbo-waddle/internal/config"
)

func main() {
	c := config.NewGopherMartConfig()

	app := app.NewServerApp(c)
	app.Run()
}
