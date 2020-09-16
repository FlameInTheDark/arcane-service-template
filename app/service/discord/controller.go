package discord

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Send simple plain text message
func (s *Service) SendMessage(channel, content string) error {
	_, err := s.session.ChannelMessageSend(channel, content)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("discord-channel", channel))
		return err
	}
	return nil
}

// Send specified discord message, like plain text, embedded or binary
func (s *Service) SendComplex(channel string, msg *discordgo.MessageSend) error {
	_, err := s.session.ChannelMessageSendComplex(channel, msg)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("discord-channel", channel))
		return err
	}
	return nil
}

// Send simple embedded message with title and plain text content
func (s *Service) SendSimpleEmbed(channel, title, content string) error {
	msg := NewEmbed("").Field(title, content, false).Color(0x00ff00).GetMessageSend()
	return s.SendComplex(channel, msg)
}
