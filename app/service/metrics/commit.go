package metrics

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"time"
)

func (s *Service) NatsMessage(subject string) {
	s.RLock()
	defer s.RUnlock()
	p := influxdb2.NewPointWithMeasurement("nats").AddField("subject", subject).SetTime(time.Now())
	s.write.WritePoint(p)
}

func (s *Service) Startup(app string) {
	s.RLock()
	defer s.RUnlock()
	p := influxdb2.NewPointWithMeasurement("startup").AddField("application", app).SetTime(time.Now())
	s.write.WritePoint(p)
}

func (s *Service) Command(id, guild string) {
	s.RLock()
	defer s.RUnlock()
	p := influxdb2.NewPointWithMeasurement("command").AddField("id", id).AddField("guild", guild).SetTime(time.Now())
	s.write.WritePoint(p)
}
