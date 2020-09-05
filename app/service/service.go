package service

import (
	"fmt"
	"github.com/FlameInTheDark/arcane-service-template/app/service/metrics"

	"github.com/FlameInTheDark/arcane-service-template/app/service/config"
	"github.com/FlameInTheDark/arcane-service-template/app/service/database"
	"github.com/FlameInTheDark/arcane-service-template/app/service/discord"
	"github.com/FlameInTheDark/arcane-service-template/app/service/etcd"
	"github.com/FlameInTheDark/arcane-service-template/app/service/log"
	"github.com/FlameInTheDark/arcane-service-template/app/service/nats"
	"go.uber.org/zap"
)

const (
	appName       = "arcane-service"
	moduleName    = "service"
	logCollection = "log"
)

var (
	zapModule      = zap.String("module", moduleName)
	zapApplication = zap.String("app", appName)
)

type Service struct {
	Logger   *zap.Logger
	Etcd     *etcd.EtcdService
	Database *database.DatabaseService
	Discord  *discord.DiscordService
	Nats     *nats.NatsService
	Metrics  *metrics.MetricsService
	Config   *config.Service
}

func New(endpoints []string, username, password string) (*Service, error) {
	logger := log.MakeLogger()

	var service Service
	service.Logger = logger
	etcdService, err := etcd.New(endpoints, username, password)
	if err != nil {
		logger.Error(err.Error(), zapModule)
		return nil, fmt.Errorf("error creating etcd session: %s", err)
	}
	etcdService.SetLogger(logger)
	service.Etcd = etcdService
	service.Config = &config.Service{}

	err = service.loadConfig()
	if err != nil {
		return nil, err
	}

	err = service.init()
	if err != nil {
		return nil, err
	}

	err = service.registerConfigWatchers()
	if err != nil {
		service.Logger.Error(err.Error(), zapModule)
		return nil, err
	}
	service.Metrics.Startup(appName)
	service.Logger.Info("Application started", zap.String("action", "launch"))
	return &service, nil
}

func (s *Service) loadConfig() error {
	err := s.Etcd.GetOneJSON(config.EtcdEnvironment, &s.Config.Environment)
	err = s.Etcd.GetOneJSON(config.EtcdDatabase, &s.Config.Database)
	err = s.Etcd.GetOneJSON(config.EtcdNats, &s.Config.Nats)
	err = s.Etcd.GetOneJSON(config.EtcdMetrics, &s.Config.Metrics)
	if err != nil {
		return fmt.Errorf("error load config: %s", err)
	}
	return nil
}

func (s *Service) init() error {
	databaseService, err := database.New(s.Config.Database.GenerateConnString(), s.Config.Environment.Database)
	if err != nil {
		s.Logger.Error(err.Error(), zapModule)
		return fmt.Errorf("init database service error: %s", err)
	}

	logger := log.MakeLoggerWriter(databaseService.MakeWriter(logCollection))
	s.Logger = logger.With(zapApplication)
	s.Etcd.SetLogger(logger)

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

	metricsService := metrics.New(s.Config.Metrics.Endpoint, s.Config.Metrics.Token, s.Config.Metrics.Org, s.Config.Metrics.Bucket)

	s.Metrics = metricsService
	s.Nats = natsService
	s.Discord = discordService
	s.Database = databaseService
	return nil
}

func (s *Service) Close() {
	s.Logger.Info("Shutting down application", zap.String("action", "shutdown"))
	s.Etcd.Close()
	s.Nats.Close()
	s.Discord.Close()
	s.Database.Close()
	s.Metrics.Close()
	s.Logger.Sync()
}
