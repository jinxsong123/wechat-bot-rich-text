package sender

type Sender interface {
	Send([]*string) error
}
