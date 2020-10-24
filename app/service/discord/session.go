package discord

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var (
	zapModule = zap.String("module", "discord")
)

type Service struct {
	sync.RWMutex
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

func (s *Service) GetMyUsername() (string, error) {
	user, err := s.session.User("@me")
	if err != nil {
		s.logger.Warn(err.Error())
		return "", err
	}
	return user.Username, nil
}

func (s *Service) Close() {
	err := s.session.Close()
	if err != nil {
		s.logger.Warn(err.Error())
	}
}

func (s *Service) SetLogger(logger *zap.Logger) {
	s.logger = logger.With(zapModule)
}

func (s *Service) Reload(token string) error {
	s.Lock()
	defer s.Unlock()
	s.Close()
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("error during session creation: %s", err)
	}

	err = sess.Open()
	if err != nil {
		return err
	}
	s.session = sess
	return nil
}
