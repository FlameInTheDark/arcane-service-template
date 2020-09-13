package database

import (
	"encoding/json"
	"io"

	"github.com/globalsign/mgo"
	"go.uber.org/zap"
)

const (
	moduleName = "database"
)

var (
	zapModule = zap.String("module", moduleName)
)

type Service struct {
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
	sess       *mgo.Session
	database   string
	collection string
}

func (dw DBWriter) Write(p []byte) (n int, err error) {
	var m map[string]interface{}
	err = json.Unmarshal(p, &m)
	if err != nil {
		return
	}
	c := dw.sess.DB(dw.database).C(dw.collection)
	err = c.Insert(m)
	if err != nil {
		return
	}
	return len(p), nil
}

func (s *Service) MakeWriter(collection string) io.Writer {
	return DBWriter{
		sess:       s.session,
		database:   s.database,
		collection: collection,
	}
}

func (s *Service) SetDatabase(database string) {
	s.database = database
}

func (s *Service) db() *mgo.Database {
	return s.session.DB(s.database)
}
