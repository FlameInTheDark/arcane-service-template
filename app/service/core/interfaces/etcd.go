package interfaces

import "go.uber.org/zap"

type EtcdService interface {
	Close()
	SetLogger(logger *zap.Logger)

	GetOneRaw(key string) ([]byte, error)
	GetAllRaw(key string) (map[string][]byte, error)
	GetOneJSON(key string, v interface{}) error
	AddWatcher(key string, handler func(key, value string, version int64)) error
	RemoveWatcher(key string)
}
