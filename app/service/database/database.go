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

type DatabaseService struct {
	session  *mgo.Session
	database string
	logger   *zap.Logger
}

func New(conn, database string) (*DatabaseService, error) {
	sess, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	return &DatabaseService{session: sess, database: database}, nil
}

func (s *DatabaseService) SetLogger(logger *zap.Logger) {
	s.logger = logger.With(zapModule)
}

func (s *DatabaseService) Close() {
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

func (s *DatabaseService) MakeWriter(collection string) io.Writer {
	return DBWriter{
		sess:       s.session,
		database:   s.database,
		collection: collection,
	}
}

func (s *DatabaseService) SetDatabase(database string) {
	s.database = database
}
