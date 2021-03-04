package models

import (
	"encoding/json"
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/sErr"
	"github.com/zqjzqj/instantCustomer/services"
	"time"
)

const VisitorsExtendsJsonKeyMessage = "message"

type Visitors struct {
	FieldsModel
	MchId uint64 `gorm:"index:idx_mchId;not null;default:0;comment:对应的商户id"`//商户id
	MAId uint64  `gorm:"index:idx_mchId;default:0;comment:对应的商户客服id"`//对应的商户客服id
	ConnId string `gorm:"size:200;not null;index:idx_connId,unique;comment:ws的连接id"`//当前连接id
	Ip string `gorm:"size:128;comment:支持ipv6的长度"`
	Country string `gorm:"size:100;default:''"` //国家
	Region string `gorm:"size:100;default:''"`//省份
	City string `gorm:"size:100;default:''"`//城市
	OnlineStatus uint8
	LastActiveTime time.Time `gorm:"index:idx_lat;comment:最近活跃时间"`
	FieldsExtendsJsonLongType

	//当前链接
	Conn *websocket.NSConn `gorm:"-"`
	Room *websocket.Room `gorm:"-"`
	Messages []*services.WsMessage `gorm:"-"`
	CurrentMa *MchAccount `gorm:"-"`
}

func (v *Visitors) TableName() string {
	return "visitors"
}

func (v *Visitors) SendMessage(data, _type string) (*services.WsMessage, error) {
	if v.MAId == 0 {
		return nil, sErr.New("未接入到客服")
	}
	msg := &services.WsMessage{
		Id: time.Now().UnixNano(),
		VisitorId:  v.ID,
		MAccountId: v.MAId,
		MchId:      v.MchId,
		Data:       data,
		Type:       _type,
		CreatedAt:  time.Now(),
		Status: services.WsMessageStatusWaitSend,
	}
	b, _ := json.Marshal(msg)
	if v.CurrentMa.Conn.Emit(services.EventOnChat, b) {
		msg.Status = services.WsMessageStatusFail
	} else {
		msg.Status = services.WsMessageStatusSend
	}
	v.Messages = append([]*services.WsMessage{
		msg,
	}, v.Messages...)
	return msg, nil
}

func (v *Visitors) ReadMessageFromStore() {
	r := v.GetExtendsJson(VisitorsExtendsJsonKeyMessage).Array()
	for _, m := range r {
		msg := &services.WsMessage{}
		err := json.Unmarshal([]byte(m.Raw), msg)
		if err == nil {
			v.Messages = append(v.Messages, msg)
		}
	}
}

func (v *Visitors) SaveMessage() {
	v.ReadMessageFromStore()
	v.SetExtendsJson(VisitorsExtendsJsonKeyMessage, v.Messages)
	config.GetDbDefault().Model(v).Save(v)
	v.Messages = v.Messages[:0]
	return
}
