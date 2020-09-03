package etcd

import (
	"fmt"
	"sync"

	"go.etcd.io/etcd/clientv3"
)

type WatcherService struct {
	sync.Mutex
	channels map[string]*clientv3.WatchChan
	close    map[string]*chan interface{}
}

func NewWatcherService() *WatcherService {
	return &WatcherService{
		channels: make(map[string]*clientv3.WatchChan),
		close:    make(map[string]*chan interface{}),
	}
}

func (w *WatcherService) saveChannel(key string, c *clientv3.WatchChan) error {
	if ok := w.channels[key]; ok != nil {
		return fmt.Errorf("watcher channel already exists")
	}
	w.channels[key] = c
	return nil
}

func (w *WatcherService) saveClose(key string, c *chan interface{}) error {
	if ok := w.close[key]; ok != nil {
		return fmt.Errorf("close channel already exists")
	}
	w.close[key] = c
	return nil
}

func (w *WatcherService) AddWatcher(key string, watch *clientv3.WatchChan, close *chan interface{}) error {
	w.Lock()
	defer w.Unlock()
	err := w.saveChannel(key, watch)
	if err != nil {
		return fmt.Errorf("add watcher error: %s", err)
	}
	err = w.saveClose(key, close)
	if err != nil {
		return fmt.Errorf("add watcher error: %s", err)
	}
	return nil
}

func (w *WatcherService) StopWatcher(key string) {
	if ok := w.close[key]; ok != nil {
		select {
		case *w.close[key] <- true:
		default:
			return
		}
	}
}

func (w *WatcherService) RemoveWatcher(key string) {
	if ok := w.channels[key]; ok != nil {
		delete(w.channels, key)
	}
	if ok := w.close[key]; ok != nil {
		w.StopWatcher(key)
		close(*w.close[key])
		delete(w.close, key)
	}
}
