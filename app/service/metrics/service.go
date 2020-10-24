package metrics

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"sync"
)

type Service struct {
	sync.RWMutex
	client influxdb2.Client
	write  api.WriteAPI
	bucket string
	org    string
}

func New(endpoint, token, org, bucket string) *Service {
	client := influxdb2.NewClient(endpoint, token)
	writeApi := client.WriteAPI(org, bucket)
	return &Service{client: client, write: writeApi, org: org, bucket: bucket}
}

func (s *Service) Close() {
	s.write.Close()
	s.client.Close()
}

func (s *Service) Reload(endpoint, token, org, bucket string) {
	s.Lock()
	s.Unlock()
	client := influxdb2.NewClient(endpoint, token)
	writeApi := client.WriteAPI(org, bucket)
	s.client = client
	s.write = writeApi
	s.org = org
	s.bucket = bucket
}
