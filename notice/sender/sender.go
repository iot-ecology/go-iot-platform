package sender

type Sender interface {
	Send(msg string) error
}
