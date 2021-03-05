package models

import (
	"encoding/json"
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/services"
	"sync"
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
	Messages services.WsMessageSlice `gorm:"-"`
	CurrentMa *MchAccount `gorm:"-"`
	Mu sync.Mutex `gorm:"-"`
}

func (v *Visitors) TableName() string {
	return "visitors"
}

func (v *Visitors) IsConnOnline() bool {
	cm := v.Conn.Conn.Server().GetConnectionsByNamespace(services.WsNamespaceVisitor)
	_, ok := cm[v.Conn.Conn.ID()]
	if ok {

		return true
	}
	return false
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

func (v *Visitors) SaveMessageLocked() {
	v.ReadMessageFromStore()
	v.SetExtendsJson(VisitorsExtendsJsonKeyMessage, v.Messages)
	config.GetDbDefault().Model(v).Save(v)
	v.Messages = v.Messages[:0]
	return
}
