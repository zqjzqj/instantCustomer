package ws

import "encoding/json"

const (
	MessageTypeString = 0
)

type MessageData struct {
	MessageType int
	Data interface{}
}

func (msgData *MessageData) ToJson() ([]byte, error) {
	return json.Marshal(msgData)
}

func (msgData *MessageData) ToMessage(from, to string) *Message {
	return &Message{
		From:        from,
		To:          to,
		MessageType: msgData.MessageType,
		Data:        msgData.Data,
	}
}

type Message struct {
	From        string
	To          string
	MessageType int
	Data        interface{}
}

func NewMessageByJsonBytes(b []byte, from, to string) (*Message, error) {
	msgData := &MessageData{}
	err := json.Unmarshal(b, msgData)
	if err != nil {
		return nil, err
	}
	return msgData.ToMessage(from, to), nil
}
