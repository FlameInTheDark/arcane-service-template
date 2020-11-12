package controller

import (
	"fmt"

	"github.com/FlameInTheDark/arcane-service-template/app/service/core/model"
	"github.com/FlameInTheDark/arcane-service-template/app/service/discord/embed"
	cModel "github.com/FlameInTheDark/arcane-service-template/model"
)

func (w *Worker) helpCommand() {
	_ = w.service.Nats.SubscribeQueue(natsWorker, "workers", func(c *cModel.Command) {
		var cmd model.GuildCommand
		var accepted bool
		err := w.service.Database.GetCommand(command, c.GuildID, &cmd)
		if err != nil {
			accepted = true
		} else {
			accepted = cmd.Active
		}
		if accepted {
			msg := embed.NewEmbed("").
				Field("Help", "Help message!", false).
				Footer(fmt.Sprintf("Requested by %s", c.Username)).
				Color(0x00ff00).
				GetMessageSend()
			err := w.service.Discord.SendComplex(c.ChannelID, msg)
			if err == nil {
				w.service.Metrics.Command(command, c.GuildID)
				_ = w.service.Database.SetUsage(command, c.UserID, c.GuildID)
			}
		}
	})
}
