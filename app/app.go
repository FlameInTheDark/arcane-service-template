package app

import (
	"fmt"
	"strings"

	"github.com/FlameInTheDark/arcane-service-template/app/service"
)

type Controller interface {
	Init(service *service.Services)
	RegisterWorkers()
}

type Application struct {
	service    *service.Services
	controller Controller
}

func New(endpoints, username, password, appName string) (*Application, error) {
	newService, err := service.New(parseEndpoints(endpoints), username, password, appName)
	if err != nil {
		return nil, fmt.Errorf("creating service error: %s", err)
	}
	return &Application{service: newService}, nil
}

func (app *Application) Start(c Controller) {
	app.controller = c
	app.controller.Init(app.service)
	app.controller.RegisterWorkers()
}

func (app *Application) Close() {
	app.service.Close()
}

func parseEndpoints(raw string) []string {
	trim := strings.ReplaceAll(raw, " ", "")
	return strings.Split(trim, ",")
}
