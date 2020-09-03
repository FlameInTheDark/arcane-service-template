package app

import (
	"errors"
	"fmt"
	"github.com/FlameInTheDark/arcane-service-template/app/service"
	"os"
	"strings"
)

var (
	etcdEndpoints = os.Getenv("ETCD_ENDPOINTS")
	etcdUsername  = os.Getenv("ETCD_USERNAME")
	etcdPassword  = os.Getenv("ETCD_PASSWORD")
)

type Application struct {
	Service *service.Service
}

func New() (*Application, error) {
	err := checkEnvironment()
	if err != nil {
		return nil, fmt.Errorf("[Application] starting application error: %s", err)
	}
	newService, err := service.New(parseEndpoints(), etcdUsername, etcdPassword)
	if err != nil {
		return nil, fmt.Errorf("[Application] creating service error: %s", err)
	}
	return &Application{Service: newService}, nil
}

func checkEnvironment() error {
	if etcdEndpoints == "" {
		return errors.New("etcd endpoints not set")
	}
	if etcdUsername == "" {
		return errors.New("etcd username not set")
	}
	if etcdPassword == "" {
		return errors.New("etcd password not set")
	}
	return nil
}

func parseEndpoints() []string {
	trim := strings.ReplaceAll(etcdEndpoints, " ", "")
	return strings.Split(trim, ",")
}
