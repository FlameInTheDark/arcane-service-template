package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlameInTheDark/arcane-service-template/app"
)

var (
	etcdEndpoints = os.Getenv("ETCD_ENDPOINTS")
	etcdUsername  = os.Getenv("ETCD_USERNAME")
	etcdPassword  = os.Getenv("ETCD_PASSWORD")
)

func main() {
	err := checkEnvironment()
	if err != nil {
		log.Fatal(err.Error())
	}
	application, err := app.New(etcdEndpoints, etcdUsername, etcdPassword)
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

func checkEnvironment() error {
	if etcdEndpoints == "" {
		return errors.New("endpoints not set")
	}
	if etcdUsername == "" {
		return errors.New("username not set")
	}
	if etcdPassword == "" {
		return errors.New("password not set")
	}
	return nil
}