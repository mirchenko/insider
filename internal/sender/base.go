package sender

type Sender interface {
	IsStarted() bool
	Start() error
	Stop()
}
