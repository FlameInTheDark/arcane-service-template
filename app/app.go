package app

import (
	"fmt"
	"strings"

	"github.com/FlameInTheDark/arcane-service-template/app/service"
)

type Application struct {
	Service *service.Service
}

func New(endpoints, username, password string) (*Application, error) {
	newService, err := service.New(parseEndpoints(endpoints), username, password)
	if err != nil {
		return nil, fmt.Errorf("creating service error: %s", err)
	}
	return &Application{Service: newService}, nil
}

func parseEndpoints(raw string) []string {
	trim := strings.ReplaceAll(raw, " ", "")
	return strings.Split(trim, ",")
}
