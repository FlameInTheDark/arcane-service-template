// Controller example package.
package controller

import (
	"github.com/FlameInTheDark/arcane-service-template/app/service"
	"github.com/FlameInTheDark/arcane-service-template/model"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Worker struct {
	service *service.Services
}

// Create new controller
func New() *Worker {
	return &Worker{}
}

// Initiate controller with application service
func (w *Worker) Init(service *service.Services) {
	w.service = service
	w.registerCommand()
}

// Register workers and commands
func (w *Worker) RegisterWorkers() {
	w.pingWorker()
	w.commandPingWorker()

	w.helpCommand()
}

func (w *Worker) pingWorker() {
	_ = w.service.Nats.Subscribe(natsPing, func(c *model.Command) {
		_ = w.service.Nats.Publish(natsPingResponse, model.RegisterCommand{Name: command, Worker: natsWorker})
	})
}

func (w *Worker) commandPingWorker() {
	_ = w.service.Nats.SubscribeQueue(natsWorkerPing, "worker", func(msg nats.Msg) {
		err := msg.Respond([]byte("pong"))
		if err != nil {
			w.service.Logger.Warn(err.Error(), zap.String("nats-subject", natsWorkerPing))
			return
		}
	})
}

func (w *Worker) registerCommand() {
	cmd := model.RegisterCommand{
		Name:   command,
		Worker: natsWorker,
	}
	_ = w.service.Nats.Publish(natsRegisterCommand, cmd)
}
