package main

import (
	"context"
	"fmt"
	"log"

	"github.com/itaraxa/turbo-waddle/internal/app"
	"github.com/itaraxa/turbo-waddle/internal/config"
	"github.com/itaraxa/turbo-waddle/internal/version"
)

func main() {
	c := config.NewGopherMartConfig()
	err := c.Config()
	if err != nil {
		log.Fatal(err)
	}
	if c.ShowVersion {
		fmt.Printf("App version: %s\n\rDatabase schema version: %d", version.ServerApp, version.Database)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.NewServerApp(ctx, c)
	app.Run()
}
