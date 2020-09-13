package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var (
	zapModule = zap.String("module", "discord")
)

type Service struct {
	session *discordgo.Session
	logger  *zap.Logger
}

func New(token string, log *zap.Logger) (*Service, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error while try to create session: %s", err)
	}
	err = sess.Open()
	if err != nil {
		return nil, err
	}
	return &Service{
		session: sess,
		logger:  log.With(zapModule),
	}, nil
}

func (d *Service) GetMyUsername() (string, error) {
	user, err := d.session.User("@me")
	if err != nil {
		d.logger.Warn(err.Error())
		return "", err
	}
	return user.Username, nil
}

func (d *Service) Close() {
	err := d.session.Close()
	if err != nil {
		d.logger.Warn(err.Error())
	}
}

func (d *Service) SetLogger(logger *zap.Logger) {
	d.logger = logger.With(zapModule)
}
