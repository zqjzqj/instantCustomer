package ws

import (
	"context"
	"encoding/json"
	"github.com/zqjzqj/instantCustomer/logs"
	"github.com/zqjzqj/instantCustomer/sErr"
)

const (
	MessageTypeString = 0
	MessageTypeClosed = 10
)

type Message struct {
	From        string      `json:"from"`
	To          string      `json:"to"`
	MessageType int         `json:"message_type"`
	Mt          int         `json:"-"`
	Data        interface{} `json:"data"`
}

func NewMessageString(mt int, from, to, data string) *Message {
	return &Message{
		From:        from,
		To:          to,
		MessageType: MessageTypeString,
		Mt:          mt,
		Data:        data,
	}
}

func NewMessageByJsonBytes(b []byte, mt int, from, to string) (*Message, error) {
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
	if msg.Data == nil {
		return nil, sErr.New("data is invalid")
	}
	msg.Mt = mt
	return msg, nil
}

func (msg *Message) ToMsgBytes() []byte {
	switch msg.MessageType {
	case MessageTypeString:
		return []byte(msg.Data.(string))
	}
	return []byte{}
}

func HandleMsgForwardToClient(msgChan chan *Message, ctx context.Context, afterFn func(msg *Message)) {
	for {
		select {
		case msg := <-msgChan:
			//转发至客户端接收
			wsc, ok := FindWsConn(msg.To)
			if !ok {
				logs.PrintlnWarning("invalid receiver ", msg.To, " msg: ", msg.Data)
				continue
			}
			err := wsc.WriteMsg(msg)
			if err != nil {
				logs.PrintErr("send msg err ", err)
			} else {
				logs.PrintlnInfo("send msg ok ", msg.Data)
			}
			if afterFn != nil {
				afterFn(msg)
			}
		case <-ctx.Done():
			logs.PrintlnInfo("exit msg forward....")
			return
		}
	}
}
