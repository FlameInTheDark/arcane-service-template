package nats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	zapModule = zap.String("module", "nats")
)

type NatsService struct {
	conn   *nats.Conn
	econn  *nats.EncodedConn
	logger *zap.Logger
}

func New(endpoints string, log *zap.Logger) (*NatsService, error) {
	natsLogger := log.With(zapModule)
	conn, err := nats.Connect(endpoints)
	if err != nil {
		return nil, err
	}

	econn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return &NatsService{
		conn:   conn,
		econn:  econn,
		logger: natsLogger,
	}, nil
}

func (n *NatsService) Close() {
	n.econn.Close()
	n.conn.Close()
}
