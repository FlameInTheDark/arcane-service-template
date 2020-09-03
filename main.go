package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/FlameInTheDark/arcane-service-template/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		return
	}

	go gracefulShutdown(application.Service.Close)
	lock := make(chan struct{})
	<-lock
}

func gracefulShutdown(close func()) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		close()
		os.Exit(0)
	}()
}
