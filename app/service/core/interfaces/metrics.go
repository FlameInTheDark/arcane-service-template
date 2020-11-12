package interfaces

import "io"

type MetricsService interface {
	Close()
	Reload(endpoint, token, org, bucket string)
	MakeWriter() io.Writer

	Startup(app string)
	NatsMessage(subject string)
	Command(id, guild string)
	CommandTotal(id string)
}
