package controller

import (
	"fmt"
	model "github.com/FlameInTheDark/arcane-service-template/app/model/database"
	natsModel "github.com/FlameInTheDark/arcane-service-template/app/model/nats"
)

func (w *Worker) pingCommand() {
	_ = w.service.Nats.SubscribeQueue(natsWorker, "workers", func(c *natsModel.Command) {
		var cmd model.GuildCommand
		var accepted bool
		err := w.service.Database.GetCommand(command, c.GuildID, &cmd)
		if err != nil {
			accepted = true
		} else {
			accepted = cmd.Active
		}
		if accepted {
			_ = w.service.Discord.SendMessage(c.ChannelID, fmt.Sprintf("<@%s> Help message!", c.UserID))
		}
	})
}
