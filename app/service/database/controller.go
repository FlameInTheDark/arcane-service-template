package database

import (
	model "github.com/FlameInTheDark/arcane-service-template/app/model/database"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/zap"
	"time"
)

func (s *Service) GetCommand(command, guild string, result interface{}) error {
	s.RLock()
	defer s.RUnlock()
	err := s.db().C("command").Find(bson.M{"id": command, "guildid": guild}).One(result)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("command", command), zap.String("guild", guild))
		_ = s.SetNewCommand(command, guild)
	}
	return err
}

func (s *Service) SetNewCommand(command, guild string) error {
	s.RLock()
	defer s.RUnlock()
	err := s.db().C("command").Insert(model.GuildCommand{
		GuildID: guild,
		ID:      command,
		Active:  true,
	})
	if err != nil {
		s.logger.Error(err.Error(), zap.String("command", command), zap.String("guild", guild))
	}
	return err
}

func (s *Service) SetUsage(command, user, guild string) error {
	s.RLock()
	defer s.RUnlock()
	err := s.db().C("action").Insert(model.GuildActions{
		UserID:    user,
		GuildID:   guild,
		Command:   command,
		Timestamp: time.Now(),
	})
	if err != nil {
		s.logger.Error(err.Error(), zap.String("command", command), zap.String("guild", guild))
	}
	return err
}
