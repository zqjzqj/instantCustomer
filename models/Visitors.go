package models

import (
	"time"
)

const VisitorsExtendsJsonKeyMessage = "message"

type Visitors struct {
	FieldsModel
	MchId          uint64 `gorm:"index:idx_mchId;not null;default:0;comment:对应的商户id"`        //商户id
	MAId           uint64 `gorm:"index:idx_mchId;default:0;comment:对应的商户客服id"`               //对应的商户客服id
	ConnId         string `gorm:"size:200;not null;index:idx_connId,unique;comment:ws的连接id"` //当前连接id
	Ip             string `gorm:"size:128;comment:支持ipv6的长度"`
	Country        string `gorm:"size:100;default:''"` //国家
	Region         string `gorm:"size:100;default:''"` //省份
	City           string `gorm:"size:100;default:''"` //城市
	OnlineStatus   uint8
	LastActiveTime time.Time `gorm:"index:idx_lat;comment:最近活跃时间"`
	FieldsExtendsJsonLongType
}

func (v *Visitors) TableName() string {
	return "visitors"
}
