package discord

import (
	"fmt"
	"github.com/FlameInTheDark/arcane-service-template/app/service/discord/embed"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Send simple plain text message
func (s *Service) SendMessage(channel, content string) error {
	s.RLock()
	defer s.RUnlock()
	_, err := s.session.ChannelMessageSend(channel, content)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("discord-channel", channel))
		return err
	}
	return nil
}

// Send specified discord message, like plain text, embedded or binary
func (s *Service) SendComplex(channel string, msg *discordgo.MessageSend) error {
	s.RLock()
	defer s.RUnlock()
	_, err := s.session.ChannelMessageSendComplex(channel, msg)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("discord-channel", channel))
		return err
	}
	return nil
}

// Send simple embedded message with title and plain text content
func (s *Service) SendSimpleEmbed(channel, title, content string) error {
	s.RLock()
	defer s.RUnlock()
	msg := embed.NewEmbed("").Field(title, content, false).Color(0x00ff00).GetMessageSend()
	return s.SendComplex(channel, msg)
}

// Send simple embedded error message with title, field title, content and username in footer
func (s *Service) SendErrorMessage(channel, title, field, content, username string) error {
	s.RLock()
	defer s.RUnlock()
	msg := embed.NewEmbed(title).Field(field, content, false).Color(0xff0000).Footer(fmt.Sprintf("Requested by %s", username)).GetMessageSend()
	return s.SendComplex(channel, msg)
}

// Send simple embedded error message with title, field title, content and username in footer
func (s *Service) SendWarningMessage(channel, title, field, content, username string) error {
	s.RLock()
	defer s.RUnlock()
	msg := embed.NewEmbed(title).Field(field, content, false).Color(0xffff00).Footer(fmt.Sprintf("Requested by %s", username)).GetMessageSend()
	return s.SendComplex(channel, msg)
}
