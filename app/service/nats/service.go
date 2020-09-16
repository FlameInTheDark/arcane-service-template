package nats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	zapModule = zap.String("module", "nats")
)

type Service struct {
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
func (n *Service) Close() {
	for _, v := range n.subscriptions {
		if v != nil {
			err := v.Unsubscribe()
			if err != nil {
				n.logger.Warn(err.Error(), zap.String("nats-subject", v.Subject), zap.String("nats-queue", v.Queue))
			}
		}
	}
	n.econn.Close()
	n.conn.Close()
}

func (n *Service) SubscribeQueue(subject, queue string, handler nats.Handler) error {
	sub, err := n.econn.QueueSubscribe(subject, queue, handler)
	if err != nil {
		n.logger.Warn(err.Error(), zap.String("nats-subject", subject), zap.String("nats-queue", queue))
		return err
	}
	n.subscriptions = append(n.subscriptions, sub)
	return nil
}

func (n *Service) Subscribe(subject string, handler nats.Handler) error {
	sub, err := n.econn.Subscribe(subject, handler)
	if err != nil {
		n.logger.Warn(err.Error(), zap.String("nats-subject", subject))
		return err
	}
	n.subscriptions = append(n.subscriptions, sub)
	return nil
}

func (n *Service) Publish(subject string, v interface{}) error {
	err := n.econn.Publish(subject, v)
	if err != nil {
		n.logger.Warn(err.Error(), zap.String("nats-subject", subject), zap.Reflect("nats-value", v))
		return err
	}
	return nil
}
