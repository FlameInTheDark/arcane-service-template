package metrics

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"time"
)

func (m *MetricsService) NatsMessage(subject string) {
	p := influxdb2.NewPointWithMeasurement("nats").AddField("subject", subject).SetTime(time.Now())
	m.write.WritePoint(p)
}

func (m *MetricsService) Startup(app string) {
	p := influxdb2.NewPointWithMeasurement("startup").AddField("application", app).SetTime(time.Now())
	m.write.WritePoint(p)
}