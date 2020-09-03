package service

import (
	"fmt"
	"github.com/FlameInTheDark/arcane-service-template/app/service/config"
	"github.com/FlameInTheDark/arcane-service-template/app/service/database"
	"github.com/FlameInTheDark/arcane-service-template/app/service/discord"
	"github.com/FlameInTheDark/arcane-service-template/app/service/etcd"
	"github.com/FlameInTheDark/arcane-service-template/app/service/nats"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
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
	Config   *config.ConfigService
}

func New(endpoints []string, username, password string) (*Service, error) {
	logger, err := makeLogger()
	if err != nil {
		return nil, err
	}

	var service Service

	etcdService, err := etcd.New(endpoints, username, password)
	if err != nil {
		logger.Error(err.Error(), zapModule)
		return nil, fmt.Errorf("error creating etcd session: %s", err)
	}

	service.Etcd = etcdService
	service.Config = &config.ConfigService{}

	err = service.LoadConfig()
	if err != nil {
		return nil, err
	}

	err = service.Init()
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *Service) LoadConfig() error {
	err := s.Etcd.GetOneJSON(config.EtcdEnvironment, &s.Config.Environment)
	err = s.Etcd.GetOneJSON(config.EtcdDatabase, &s.Config.Database)
	err = s.Etcd.GetOneJSON(config.EtcdNats, &s.Config.Nats)
	if err != nil {
		s.Logger.Warn(err.Error(), zapModule)
		return fmt.Errorf("error load config: %s", err)
	}
	return nil
}

func (s *Service) Init() error {
	databaseService, err := database.New(s.Config.Database.GenerateConnString(), s.Config.Environment.Database)
	if err != nil {
		s.Logger.Error(err.Error(), zapModule)
		return fmt.Errorf("init database service error: %s", err)
	}

	lg := makeDBLogger(databaseService.MakeWriter(logCollection))
	s.Logger = lg.With(zapApplication)
	s.Etcd.SetLogger(lg)

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

	s.Nats = natsService
	s.Discord = discordService
	s.Database = databaseService
	return nil
}

func (s *Service) Close() {
	s.Etcd.Close()
	s.Nats.Close()
	s.Discord.Close()
	s.Database.Close()
	s.Logger.Sync()
}

func makeLogger() (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	return cfg.Build()
}

//TODO: https://stackoverflow.com/questions/40396499/go-create-io-writer-inteface-for-logging-to-mongodb-database
func makeDBLogger(mw io.Writer) *zap.Logger {
	cfg := zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(mw)), zapcore.DebugLevel))
	return logger
}
