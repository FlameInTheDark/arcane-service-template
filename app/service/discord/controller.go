package discord

import "go.uber.org/zap"

func (s *Service) SendMessage(channel, content string) error {
	_, err := s.session.ChannelMessageSend(channel, content)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("discord-channel", channel))
		return err
	}
	return nil
}
