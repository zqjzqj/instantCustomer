package models

import (
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/sErr"
	"github.com/zqjzqj/instantCustomer/ws"
	"gorm.io/gorm"
	"time"
)

const (
	MerchantVersionTrial    = 0
	MerchantVersionStandard = 1
	MerchantVersionSenior   = 2

	MerchantDefaultMaxSeats = 1
	MerchantTrialMaxSeats   = 99
)

type Merchant struct {
	FieldsModel
	Name          string    `gorm:"size:150;index:idx_name;not null;comment:商户名称"`
	Contacts      string    `gorm:"size:50;not null;comment:联系人"`
	Phone         string    `gorm:"size:20;index:idx_phone;not null;comment:联系人电话"`
	MaxSeats      uint      `gorm:"type:int(11) unsigned;default:1;comment:最大坐席数(可以同时在线的成员账号数)"`
	Expired       time.Time `gorm:"index:idx_expired;autoCreateTime;comment:到期时间"`
	Version       uint8     `gorm:"default:0;comment:版本：0体验版 1标准版 2高级版...."`
	Status        uint8     `gorm:"default:1;comment:商户状态 1正常 0禁用"`
	LastLoginTime time.Time `gorm:"comment:最近一次登陆时间"`
	FieldsExtendsJsonType
}

func (m *Merchant) TableName() string {
	return "merchant"
}

//创建商户
func NewCreateMerchant(phone, name, password string) (*MchAccount, error) {
	now := time.Now()
	mch := &Merchant{
		Name:          name,
		Contacts:      "管理员",
		Phone:         phone,
		MaxSeats:      MerchantTrialMaxSeats,
		Expired:       now.AddDate(0, 0, 7), //体验7天
		Version:       MerchantVersionTrial, //体验版
		Status:        global.IsYes,         //商户状态正常
		LastLoginTime: now,
	}
	ma := &MchAccount{}
	db := global.GetDb()
	db.Where("phone = ?", mch.Phone).Find(ma)
	if ma.ID > 0 {
		return nil, sErr.New("改手机号已被注册")
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(mch).Error
		if err != nil {
			return err
		}
		ma, err = CreateNewAccount(mch, name, phone, password, MchAccountRoleAdmin, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ma, nil
}

func AllocationMchAccount(mchId string) (*MchAccount, error) {
RE:
	ma := &MchAccount{}
	if global.GetDb().Where("mch_id = ? and session_num < max_session", mchId).
		Where("online_status", global.OnlineStatusYes).
		Limit(1).Order("updated_at desc").Find(ma).RowsAffected == 0 {
		return nil, sErr.New("merchant reach the upper limit")
	}
	if _, ok := ws.FindWsConn(ma.ConnId.String); !ok {
		global.GetDb().Table(ma.TableName()).Where("id = ?", ma.ID).Update("online_status", global.OnlineStatusLeave)
		goto RE
	}
	return ma, nil
}
