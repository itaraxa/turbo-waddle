package main

import "github.com/itaraxa/turbo-waddle/internal/app"

func main() {
	app := app.NewServerApp(nil, nil, nil, nil)
	app.Run()
}
