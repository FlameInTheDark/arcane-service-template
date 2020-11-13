package interfaces

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type DiscordService interface {
	Close()
	SetLogger(logger *zap.Logger)
	Reload(token string) error

	GetMyUsername() (string, error)
	SendMessage(channel, content string) error
	SendComplex(channel string, msg *discordgo.MessageSend) error
	SendSimpleEmbed(channel, title, content string) error
	SendErrorEmbed(channel, title, field, content, username string) error
	SendWarningEmbed(channel, title, field, content, username string) error
}
