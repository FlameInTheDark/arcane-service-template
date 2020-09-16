package metrics

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

type Service struct {
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

func (m *Service) Close() {
	m.write.Close()
	m.client.Close()
}
