package metrics

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
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
	p := influxdb2.NewPointWithMeasurement("command").AddField(id, 1).AddTag("guild", guild).SetTime(time.Now())
	s.write.WritePoint(p)
	s.CommandTotal(id)
}

func (s *Service) CommandTotal(id string) {
	s.RLock()
	defer s.RUnlock()
	totalOne := influxdb2.NewPointWithMeasurement("command").AddField(id, 1).AddTag("guild", "total").SetTime(time.Now())
	s.write.WritePoint(totalOne)
	totalAll := influxdb2.NewPointWithMeasurement("command").AddField("total", 1).AddTag("guild", "total").SetTime(time.Now())
	s.write.WritePoint(totalAll)
}
