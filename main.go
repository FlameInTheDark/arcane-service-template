package main

import (
	"fmt"
	"github.com/FlameInTheDark/arcane-service-template/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	user, err := application.Service.Discord.GetMyUsername()
	if err != nil {
		fmt.Println(err)
		return
	}
	application.Service.Logger.Info(user)
	application.Service.Close()
}
