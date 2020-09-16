package controller

import (
	"fmt"
	model "github.com/FlameInTheDark/arcane-service-template/app/model/database"
	natsModel "github.com/FlameInTheDark/arcane-service-template/app/model/nats"
	"github.com/FlameInTheDark/arcane-service-template/app/service/discord"
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
			msg := discord.NewEmbed("").
				Field("Help", "Help message!", false).
				Footer(fmt.Sprintf("Requested by %s", c.Username)).
				Color(0x00ff00).
				GetMessageSend()

			err := w.service.Discord.SendComplex(c.ChannelID, msg)
			if err == nil {
				w.service.Metrics.Command(command, c.GuildID)
			}
		}
	})
}
