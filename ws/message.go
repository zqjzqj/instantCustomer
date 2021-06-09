package ws

const (
	MessageTypeString = 0
)

type Message struct {
	Form        string
	To          string
	MessageType int
	Data        interface{}
}

func NewMessageByBytes(b []byte) (*Message, error) {
	return nil, nil
}
