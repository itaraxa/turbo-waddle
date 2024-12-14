package main

import (
	"log"

	"github.com/itaraxa/turbo-waddle/internal/app"
	"github.com/itaraxa/turbo-waddle/internal/config"
)

func main() {
	c := config.NewGopherMartConfig()
	err := c.Config()
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewServerApp(c)
	app.Run()
}
