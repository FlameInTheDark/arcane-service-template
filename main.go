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

	application.Controller.Init()
	application.Controller.RegisterWorkers()

	lock := make(chan os.Signal, 1)
	signal.Notify(lock, os.Interrupt)
	signal.Notify(lock, syscall.SIGTERM)
	<-lock
	application.Service.Close()
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
