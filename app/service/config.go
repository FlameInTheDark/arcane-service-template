package service

import (
	"encoding/json"
	"github.com/FlameInTheDark/arcane-service-template/app/service/logging"

	"github.com/FlameInTheDark/arcane-service-template/app/service/config"
	"github.com/FlameInTheDark/arcane-service-template/app/service/database"
	"github.com/FlameInTheDark/arcane-service-template/app/service/discord"
	"github.com/FlameInTheDark/arcane-service-template/app/service/nats"
	"go.uber.org/zap"
)

var (
	zapActionConfigUpdate = zap.String("action", "config-update")
)

func (s *Services) registerConfigWatchers() error {
	err := s.Etcd.AddWatcher(config.EtcdNats, func(key, value string, version int64) {
		s.Config.Lock()
		err := json.Unmarshal([]byte(value), &s.Config.Nats)
		if err != nil {
			s.Logger.Error(
				err.Error(),
				zapModule,
				zapActionConfigUpdate,
				zap.String("etcd-key", config.EtcdNats))
			return
		}
		s.Config.Unlock()
		s.Nats.Close()
		endpoints := s.Config.Nats.GenerateConnString()
		natsService, err := nats.New(endpoints, s.Logger)
		if err != nil {
			s.Logger.Error(
				err.Error(),
				zapModule,
				zapActionConfigUpdate,
				zap.String("etcd-key", config.EtcdNats))
			return
		}
		s.Nats = natsService
		s.Logger.Info("Settings updated successfully", zap.String("etcd-key", config.EtcdNats))
	})
	if err != nil {
		s.Logger.Error(err.Error(), zapModule, zap.String("etcd-key", config.EtcdNats))
		return err
	}

	err = s.Etcd.AddWatcher(config.EtcdDatabase, func(key, value string, version int64) {
		s.Config.Lock()
		err := json.Unmarshal([]byte(value), &s.Config.Database)
		if err != nil {
			s.Logger.Error(
				err.Error(),
				zapModule,
				zapActionConfigUpdate,
				zap.String("etcd-key", config.EtcdDatabase))
			return
		}
		s.Config.Unlock()

		s.Database.Close()
		endpoints := s.Config.Database.GenerateConnString()
		databaseService, err := database.New(endpoints, s.Config.Environment.Database)
		if err != nil {
			s.Logger = logging.MakeLogger()
			s.Logger.Error(
				err.Error(),
				zapModule,
				zapActionConfigUpdate,
				zap.String("etcd-key", config.EtcdDatabase))
			return
		}

		databaseService.SetLogger(s.Logger)

		s.Database = databaseService
		s.Logger.Info("Settings updated successfully", zap.String("etcd-key", config.EtcdDatabase))
	})
	if err != nil {
		s.Logger.Error(err.Error(), zapModule, zap.String("etcd-key", config.EtcdDatabase))
		return err
	}

	err = s.Etcd.AddWatcher(config.EtcdEnvironment, func(key, value string, version int64) {
		s.Config.Lock()
		err := json.Unmarshal([]byte(value), &s.Config.Environment)
		if err != nil {
			s.Logger.Error(
				err.Error(),
				zapModule,
				zapActionConfigUpdate,
				zap.String("etcd-key", config.EtcdEnvironment))
			return
		}
		s.Config.Unlock()
		s.Database.SetDatabase(s.Config.Environment.Database)

		s.Discord.Close()
		discordService, err := discord.New(s.Config.Environment.Discord, s.Logger)
		if err != nil {
			s.Logger.Error(
				err.Error(),
				zapModule,
				zapActionConfigUpdate,
				zap.String("etcd-key", config.EtcdEnvironment))
			return
		}

		s.Discord = discordService
		s.Logger.Info("Settings updated successfully", zap.String("etcd-key", config.EtcdEnvironment))
	})
	if err != nil {
		s.Logger.Error(err.Error(), zapModule, zap.String("etcd-key", config.EtcdEnvironment))
		return err
	}
	return nil
}
