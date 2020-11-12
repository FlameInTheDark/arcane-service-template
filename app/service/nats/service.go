package nats

import (
	"sync"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	zapModule = zap.String("module", "nats")
)

type Service struct {
	sync.RWMutex
	conn  *nats.Conn
	econn *nats.EncodedConn

	subscriptions []*nats.Subscription

	logger *zap.Logger
}

func New(endpoints string, log *zap.Logger) (*Service, error) {
	natsLogger := log.With(zapModule)
	conn, err := nats.Connect(endpoints)
	if err != nil {
		return nil, err
	}

	econn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return &Service{
		conn:   conn,
		econn:  econn,
		logger: natsLogger,
	}, nil
}

// Close unsubscribes all subscriptions and close connection
func (s *Service) Close() {
	for _, v := range s.subscriptions {
		if v != nil {
			err := v.Unsubscribe()
			if err != nil {
				s.logger.Warn(err.Error(), zap.String("nats-subject", v.Subject), zap.String("nats-queue", v.Queue))
			}
		}
	}
	s.econn.Close()
	s.conn.Close()
}

func (s *Service) SubscribeQueue(subject, queue string, handler nats.Handler) error {
	sub, err := s.econn.QueueSubscribe(subject, queue, handler)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("nats-subject", subject), zap.String("nats-queue", queue))
		return err
	}
	s.subscriptions = append(s.subscriptions, sub)
	return nil
}

func (s *Service) Subscribe(subject string, handler nats.Handler) error {
	sub, err := s.econn.Subscribe(subject, handler)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("nats-subject", subject))
		return err
	}
	s.subscriptions = append(s.subscriptions, sub)
	return nil
}

func (s *Service) Publish(subject string, v interface{}) error {
	err := s.econn.Publish(subject, v)
	if err != nil {
		s.logger.Warn(err.Error(), zap.String("nats-subject", subject), zap.Reflect("nats-value", v))
		return err
	}
	return nil
}

func (s *Service) Reload(endpoints string) error {
	s.Lock()
	defer s.Unlock()
	conn, err := nats.Connect(endpoints)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	econn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	s.Close()
	s.conn = conn
	s.econn = econn
	return nil
}
