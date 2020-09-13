package metrics

import (
	"encoding/json"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"io"
	"time"
)

type MetricsLog struct {
	Level       string `json:"level"`
	Application string `json:"app"`
	Action      string `json:"action"`
}

type MetricsWriter struct {
	write  api.WriteAPI
	bucket string
	org    string
}

func (m *Service) MakeWriter() io.Writer {
	return MetricsWriter{
		write:  m.write,
		bucket: m.bucket,
		org:    m.org,
	}
}

func (w MetricsWriter) Write(p []byte) (n int, err error) {
	var m MetricsLog
	err = json.Unmarshal(p, &m)
	if err != nil {
		return
	}
	point := influxdb2.NewPointWithMeasurement("log").
		AddField("Application", m.Application).
		AddField("Action", m.Action).
		AddField("count", 1).
		AddField("Level", m.Level).
		SetTime(time.Now())
	w.write.WritePoint(point)
	return len(p), nil
}
