package ws

import "encoding/json"

const (
	MessageTypeString = 0
)

type Message struct {
	From        string      `json:"from"`
	To          string      `json:"to"`
	MessageType int         `json:"message_type"`
	Data        interface{} `json:"data"`
}

func NewMessageByJsonBytes(b []byte, from, to string) (*Message, error) {
	msg := &Message{}
	err := json.Unmarshal(b, msg)
	if err != nil {
		return nil, err
	}
	if from != "" {
		msg.From = from
	}
	if to != "" {
		msg.To = to
	}
	return msg, nil
}
