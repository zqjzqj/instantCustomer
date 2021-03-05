package services

import (
	"encoding/json"
	"github.com/zqjzqj/instantCustomer/models"
	"github.com/zqjzqj/instantCustomer/sErr"
	"time"
)

const(
	WsMessageStatusWaitSend = 0
	WsMessageStatusSend = 1
	WsMessageStatusFail = 2

	WsMessageSenderMchAccount = 1
	WsMessageSenderVisitor = 2

	EventOnChat = "OnChat"


	WsNamespaceMch = "mch"
	WsNamespaceVisitor = "visitor"
)

type WsMessage struct {
	Id int64
	VisitorId uint64
	MAccountId uint64
	MchId uint64
	Data string
	Sender uint8
	CreatedAt time.Time
	Status uint8
	MchAccount *models.MchAccount `json:"-"`
	Visitor *models.Visitors `json:"-"`
}

func SendWsMessage(from, to interface{}, data string) (*WsMessage, error) {
	now := time.Now()
	msg := &WsMessage{
		Id:         now.UnixNano(),
		VisitorId:  0,
		MAccountId: 0,
		MchId:      0,
		Data:       "",
		Sender:    0,
		CreatedAt:  now,
		Status:     WsMessageStatusWaitSend,
	}
	switch e := from.(type) {
	case *models.MchAccount:
		msg.Sender = WsMessageSenderMchAccount
		msg.MchId = e.MchId
		msg.MAccountId = e.ID
		msg.MchAccount = e
	case *models.Visitors:
		msg.Sender = WsMessageSenderVisitor
		msg.MchId = e.MchId
		msg.VisitorId = e.ID
		msg.Visitor = e
	default:
		return nil, sErr.New("无效的form")
	}

	//表示from是商户客服 所以这里校验和设置一下to的访客
	if msg.MAccountId > 0 {
		e, ok := to.(*models.Visitors)
		if !ok {
			return nil, sErr.New("无效的form")
		}
		msg.VisitorId = e.ID
		msg.Visitor = e
	} else {
		e, ok := to.(*models.MchAccount)
		if !ok {
			return nil, sErr.New("无效的form")
		}
		msg.MAccountId = e.ID
		msg.MchAccount = e
	}

	msg.Data = data
	if msg.Sender == WsMessageSenderMchAccount {
		if !msg.MchAccount.IsConnOnline() {
			return nil, sErr.New("商户已离线，消息发送失败")
		}
		if !msg.Visitor.IsConnOnline() {
			msg.StoreCache(true)
			return msg, nil
		}

		if msg.MchAccount.Conn.Emit(EventOnChat, msg.ToBytes()) {
			msg.Status = WsMessageStatusSend
		} else {
			msg.Status = WsMessageStatusFail
		}
	} else {
		if !msg.Visitor.IsConnOnline() {
			return nil, sErr.New("访客已离线，消息发送失败")
		}
		if !msg.MchAccount.IsConnOnline() {
			msg.StoreCache(true)
			return msg, nil
		}

		if msg.Visitor.Conn.Emit(EventOnChat, msg.ToBytes()) {
			msg.Status = WsMessageStatusSend
		} else {
			msg.Status = WsMessageStatusFail
		}
	}
	msg.StoreCache(true)
	return msg, nil
}

func (msg *WsMessage) ToBytes() []byte {
	b, _ := json.Marshal(msg)
	return b
}

func (msg *WsMessage) StoreCache(isSaveStore bool) {
	if msg.Visitor == nil {
		return
	}
	msg.Visitor.Mu.Lock()
	defer msg.Visitor.Mu.Unlock()
	msg.Visitor.Messages = append(WsMessageSlice{msg}, msg.Visitor.Messages...)
	if isSaveStore {
		//持久保存消息记录
		msg.Visitor.SaveMessageLocked()
	}
}


type WsMessageSlice []*WsMessage

func (wsmh WsMessageSlice) Len() int {
	return len(wsmh)
}

func (wsmh WsMessageSlice) Less(i, j int) bool {
	return wsmh[i].CreatedAt.After(wsmh[j].CreatedAt)
}

func (wsmh WsMessageSlice) Swap(i, j int) {
	wsmh[i], wsmh[j] = wsmh[j], wsmh[i]
}
