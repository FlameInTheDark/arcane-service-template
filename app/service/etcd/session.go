package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
)

const (
	DialTimeout = 5 * time.Second
)

var (
	zapModule = zap.String("module", "etcd")
)

type Service struct {
	session *clientv3.Client
	watcher *WatcherService
	logger  *zap.Logger
}

func New(endpoints []string, username, password string) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		Username:    username,
		Password:    password,
		DialTimeout: DialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("creating etcd session error: %s", err)
	}
	return &Service{session: cli, watcher: NewWatcherService()}, nil
}

func (e *Service) SetLogger(logger *zap.Logger) {
	e.logger = logger.With(zapModule)
}

func (e *Service) GetOneRaw(key string) ([]byte, error) {
	resp, err := e.session.Get(context.Background(), key)
	if err != nil {
		e.logger.Warn(err.Error())
		return nil, fmt.Errorf("get one value error: %s", err)
	}

	if len(resp.Kvs) == 0 {
		e.logger.Warn("no values found", zap.String("key", key))
		return nil, errors.New("no values found: " + key)
	}
	return resp.Kvs[0].Value, nil
}

func (e *Service) GetAllRaw(key string) (map[string][]byte, error) {
	resp, err := e.session.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		e.logger.Warn(err.Error())
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		e.logger.Warn("no values found")
		return nil, errors.New("no values found")
	}
	data := make(map[string][]byte)

	for _, v := range resp.Kvs {
		data[string(v.Key)] = v.Value
	}
	return data, nil
}

func (e *Service) GetOneJSON(key string, v interface{}) error {
	raw, err := e.GetOneRaw(key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, v)
	if err != nil {
		e.logger.Warn(err.Error())
		return fmt.Errorf("JSON unmarshal error: %s", err)
	}
	return nil
}

func (e *Service) AddWatcher(key string, handler func(key, value string, version int64)) error {
	wc := e.session.Watch(context.Background(), key)
	closeChan := make(chan interface{})

	err := e.watcher.AddWatcher(key, &wc, &closeChan)
	if err != nil {
		close(closeChan)
		return err
	}

	go func(key string, watcher *clientv3.WatchChan, close *chan interface{}, f func(key, value string, version int64)) {
		for {
			select {
			case change := <-*watcher:
				for _, e := range change.Events {
					if e.Kv != nil {
						f(key, string(e.Kv.Value), e.Kv.Version)
					}
				}
			case <-*close:
				return
			}
		}
	}(key, &wc, &closeChan, handler)
	return nil
}

func (e *Service) RemoveWatcher(key string) {
	e.watcher.RemoveWatcher(key)
}

func (e *Service) Close() {
	err := e.session.Close()
	if err != nil {
		e.logger.Warn(err.Error())
	}
}
