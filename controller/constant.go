package controller

const (
	command = "help"

	natsWorker          = "command." + command
	natsWorkerPing      = natsWorker + ".ping"
	natsPing            = "ping.command"
	natsPingResponse    = "ping.command.response"
	natsRegisterCommand = "register.command"
)
