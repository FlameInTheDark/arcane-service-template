package metrics

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

type MetricsService struct {
	client influxdb2.Client
	write  api.WriteAPI
	bucket string
	org    string
}

func New(endpoint, token, org, bucket string) *MetricsService {
	client := influxdb2.NewClient(endpoint, token)
	writeApi := client.WriteAPI(org, bucket)
	return &MetricsService{client: client, write: writeApi, org: org, bucket: bucket}
}

func (m *MetricsService) Close() {
	m.write.Close()
	m.client.Close()
}
