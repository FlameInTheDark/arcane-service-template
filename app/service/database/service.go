package database

import (
	"encoding/json"
	"errors"
	"github.com/globalsign/mgo"
	"go.uber.org/zap"
	"io"
	"sync"
)

const (
	moduleName = "database"
)

var (
	zapModule = zap.String("module", moduleName)
)

type Service struct {
	sync.RWMutex
	session  *mgo.Session
	database string
	logger   *zap.Logger
}

func New(conn, database string) (*Service, error) {
	sess, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	return &Service{session: sess, database: database}, nil
}

func (s *Service) SetLogger(logger *zap.Logger) {
	s.logger = logger.With(zapModule)
}

func (s *Service) Close() {
	s.session.Close()
}

type DBWriter struct {
	service    *Service
	database   string
	collection string
}

func (dw DBWriter) Write(p []byte) (n int, err error) {
	if dw.service.session == nil {
		return 0, errors.New("session is nil")
	}
	var m map[string]interface{}
	err = json.Unmarshal(p, &m)
	if err != nil {
		return
	}
	c := dw.service.session.DB(dw.database).C(dw.collection)
	err = c.Insert(m)
	if err != nil {
		return
	}
	return len(p), nil
}

func (s *Service) MakeWriter(collection string) io.Writer {
	return DBWriter{
		service:    s,
		database:   s.database,
		collection: collection,
	}
}

func (s *Service) SetDatabase(database string) {
	s.database = database
}

func (s *Service) db() *mgo.Database {
	s.RLock()
	defer s.RUnlock()
	return s.session.DB(s.database)
}

func (s *Service) Reload(conn, database string) error {
	s.Lock()
	defer s.Unlock()
	s.Close()
	sess, err := mgo.Dial(conn)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	s.session = sess
	s.database = database
	return nil
}
