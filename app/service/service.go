package service

import (
	"fmt"

	"github.com/FlameInTheDark/arcane-service-template/app/service/config"
	"github.com/FlameInTheDark/arcane-service-template/app/service/core/interfaces"
	"github.com/FlameInTheDark/arcane-service-template/app/service/database"
	"github.com/FlameInTheDark/arcane-service-template/app/service/discord"
	"github.com/FlameInTheDark/arcane-service-template/app/service/etcd"
	"github.com/FlameInTheDark/arcane-service-template/app/service/logging"
	"github.com/FlameInTheDark/arcane-service-template/app/service/metrics"
	"github.com/FlameInTheDark/arcane-service-template/app/service/nats"
	"go.uber.org/zap"
)

const (
	moduleName    = "service"
	logCollection = "log"
)

var (
	zapModule      = zap.String("module", moduleName)
	zapApplication zap.Field
)

type Services struct {
	Logger   *zap.Logger
	Etcd     interfaces.EtcdService
	Database interfaces.DatabaseService
	Discord  interfaces.DiscordService
	Nats     interfaces.NatsService
	Metrics  interfaces.MetricsService
	Config   *config.Service
}

func New(endpoints []string, username, password, appName string) (*Services, error) {
	logger := logging.MakeLogger()
	zapApplication = zap.String("app", appName)

	var services Services
	services.Logger = logger
	etcdService, err := etcd.New(endpoints, username, password)
	if err != nil {
		logger.Error(err.Error(), zapModule)
		return nil, fmt.Errorf("error creating etcd session: %s", err)
	}
	etcdService.SetLogger(logger)
	services.Etcd = etcdService
	services.Config = &config.Service{}

	err = services.loadConfig()
	if err != nil {
		return nil, err
	}

	err = services.init()
	if err != nil {
		return nil, err
	}

	err = services.registerConfigWatchers()
	if err != nil {
		services.Logger.Error(err.Error(), zapModule)
		return nil, err
	}
	services.Metrics.Startup(appName)
	services.Logger.Info("Application started", zap.String("action", "launch"))
	return &services, nil
}

func (s *Services) loadConfig() error {
	err := s.Etcd.GetOneJSON(config.EtcdEnvironment, &s.Config.Environment)
	if err != nil {
		return fmt.Errorf("error load config: %s", err)
	}

	err = s.Etcd.GetOneJSON(config.EtcdDatabase, &s.Config.Database)
	if err != nil {
		return fmt.Errorf("error load config: %s", err)
	}

	err = s.Etcd.GetOneJSON(config.EtcdNats, &s.Config.Nats)
	if err != nil {
		return fmt.Errorf("error load config: %s", err)
	}

	err = s.Etcd.GetOneJSON(config.EtcdMetrics, &s.Config.Metrics)
	if err != nil {
		return fmt.Errorf("error load config: %s", err)
	}
	return nil
}

func (s *Services) init() error {
	databaseService, err := database.New(s.Config.Database.GenerateConnString(), s.Config.Environment.Database)
	if err != nil {
		s.Logger.Error(err.Error(), zapModule)
		return fmt.Errorf("init database service error: %s", err)
	}
	metricsService := metrics.New(s.Config.Metrics.Endpoint, s.Config.Metrics.Token, s.Config.Metrics.Org, s.Config.Metrics.Bucket)

	logger := logging.MakeLoggerWriter(databaseService.MakeWriter(logCollection), metricsService.MakeWriter())
	_ = s.Logger.Sync()
	s.Logger = logger.With(zapApplication)

	s.Etcd.SetLogger(logger)
	databaseService.SetLogger(s.Logger)

	discordService, err := discord.New(s.Config.Environment.Discord, s.Logger)
	if err != nil {
		s.Logger.Error(err.Error(), zapModule)
		return fmt.Errorf("init discord service error: %s", err)
	}

	natsService, err := nats.New(s.Config.Nats.GenerateConnString(), s.Logger)
	if err != nil {
		s.Logger.Error(err.Error(), zapModule)
		return fmt.Errorf("init nats service error: %s", err)
	}

	s.Metrics = metricsService
	s.Nats = natsService
	s.Discord = discordService
	s.Database = databaseService
	return nil
}

func (s *Services) Close() {
	s.Logger.Info("Shutting down the application", zap.String("action", "shutdown"))
	defer func() { _ = s.Logger.Sync() }()
	s.Etcd.Close()
	s.Nats.Close()
	s.Discord.Close()
	s.Database.Close()
	s.Metrics.Close()
}
