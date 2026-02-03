package main

import (
	"os"
	"os/signal"
	"syscall"
	"taskservice/internal/app"
)

func main() {
	app := app.NewApp()

	go app.Run()

	sysCh := make(chan os.Signal, 1)
	signal.Notify(sysCh, syscall.SIGTERM, syscall.SIGINT)

	<-sysCh

	app.Stop()
}
