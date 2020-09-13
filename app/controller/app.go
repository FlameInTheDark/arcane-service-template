package controller

import (
	natsModel "github.com/FlameInTheDark/arcane-service-template/app/model/nats"
	"github.com/FlameInTheDark/arcane-service-template/app/service"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Worker struct {
	service *service.Service
}

func New(service *service.Service) *Worker {
	return &Worker{service: service}
}

func (w *Worker) Init() {
	w.registerCommand()
}

func (w *Worker) RegisterWorkers() {
	w.pingWorker()
	w.commandPingWorker()

	w.pingCommand()
}

func (w *Worker) pingWorker() {
	_ = w.service.Nats.Subscribe(natsPing, func(c *natsModel.Command) {
		_ = w.service.Nats.Publish(natsPingResponse, natsModel.RegisterCommand{Name: command, Worker: natsWorker})
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
	cmd := natsModel.RegisterCommand{
		Name:   command,
		Worker: natsWorker,
	}
	_ = w.service.Nats.Publish(natsRegisterCommand, cmd)
}
