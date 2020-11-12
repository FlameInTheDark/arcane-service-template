package interfaces

import (
	"io"

	"go.uber.org/zap"
)

type DatabaseService interface {
	Close()
	SetLogger(logger *zap.Logger)
	MakeWriter(collection string) io.Writer
	SetDatabase(database string)
	Reload(conn, database string) error

	GetCommand(command, guild string, result interface{}) error
	SetNewCommand(command, guild string) error
	SetUsage(command, user, guild string) error
}
