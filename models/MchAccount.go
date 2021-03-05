package models

import (
	"database/sql"
	"fmt"
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/sErr"
	"github.com/zqjzqj/instantCustomer/services"
	"gorm.io/gorm"
	"time"
)

const(
	MchAccountMaxSessionDefault = 99

	//角色常量
	MchAccountRoleAdmin = 2
	MchAccountRoleCustomer = 1
)

type MchAccount struct {
	FieldsModel
	MchId uint64 `gorm:"index:idx_mchId;not null;default:0"`
	Phone string `gorm:"size:20;comment:手机号;not null;index:idx_phone,unique"`
	Password string `gorm:"type:char(32);default:'';comment:密码md5"`
	Salt string `gorm:"type:varchar(32);default:'';comment:盐"`
	Avatar string `gorm:"type:text;comment:头像"`
	WxOpenId sql.NullString `gorm:"size:200;comment:微信openid;index:idx_wx,unique"`
	RealName string `gorm:"size:20;comment:姓名;default:''"`
	MaxSession uint `gorm:"type:int(11) unsigned;default:99;comment:最大会话数"`
	Role uint8	`gorm:"default:1;comment:1普通客服 2管理员"`
	Token  sql.NullString `gorm:"size:32;index:idx_token,unique';comment:用户移动端登陆token"`
	TokenStatus uint8 `gorm:"default:1;comment:0禁用 1正常"`
	Status uint8 `gorm:"default:1;comment:0禁用 1正常"`
	OnlineStatus uint8 `gorm:"default:0;comment:商户在线状态 1在线 0离线 2隐身"`
	LastLoginTime time.Time `gorm:"comment:最近一次登陆时间"`
	FieldsExtendsJsonType
	ConnId string `gorm:"size:200;index:idx_connId,unique;comment:ws的连接id"`//当前连接id
	Mch *Merchant `gorm:"foreignKey:Id;-"`

	SessionNum int `gorm:"default:0;comment:当前会话数"`
	Conn *websocket.NSConn `gorm:"-"`
	Room *websocket.Room `gorm:"-"`
}

func (ma *MchAccount) TableName() string {
	return "mch_account"
}

func (ma *MchAccount) IsConnOnline() bool {
	cm := ma.Conn.Conn.Server().GetConnectionsByNamespace(services.WsNamespaceMch)
	_, ok := cm[ma.Conn.Conn.ID()]
	if ok {

		return true
	}
	return false
}

func (ma *MchAccount) GetWsRoomId() string {
	return ma.GetMerchant().WsRoomId
}

func LoginMch(phone, pwd string) (*MchAccount, error) {
	ma := &MchAccount{}
	config.GetDbDefault().Where("phone = ?", phone).Find(ma)
	if ma.ID == 0 {
		return nil, sErr.New("无效的手机号")
	}
	if global.PwdPlaintext2CipherText(pwd, ma.Salt) == ma.Password {
		err := ma.LoginSuccess()
		if err != nil {
			return nil, err
		}
		return ma, nil
	}
	return nil, sErr.New("账号密码错误")
}

func CreateNewAccount(mch *Merchant, realName, phone, password string, role uint8, tx *gorm.DB) (*MchAccount, error) {
	if phone == "" || global.StrLen(phone) > 15 {
		return nil, sErr.New("无效的手机号码")
	}
	if realName == "" || global.StrLen(realName) > 20 {
		return nil, sErr.New("账户昵称不合法")
	}
	if password == "" {
		return nil, sErr.New("密码不能为空")
	}
	if role != MchAccountRoleAdmin {
		role = MchAccountRoleCustomer
	}
	ma := &MchAccount{
		MchId:                 mch.ID,
		Phone:                 phone,
		Password:              "",
		Salt:                  "",
		Avatar:                "",
		RealName:              realName,
		MaxSession:            MchAccountMaxSessionDefault,
		Role:                  role,
		Token:				   sql.NullString{
			String: "",
			Valid:  false,
		},
		TokenStatus:           global.IsNo,
		Status:                global.IsYes,
		OnlineStatus:          global.IsNo,
		WxOpenId:              sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	ma.GeneratePassword(password)
	if tx == nil {
		tx = config.GetDbDefault().DB
	}
	if err := tx.Create(ma).Error; err != nil {
		return nil, err
	}
	ma.Mch = mch
	return ma, nil
}

func (ma *MchAccount) GetMch() {
	if ma.Mch == nil {
		ma.Mch = &Merchant{}
	}
	db := config.GetDbDefault()
	db.Where("id = ?", ma.MchId).Find(ma.Mch)
}

func (ma *MchAccount) GetMerchant() *Merchant {
	if ma.Mch != nil {
		return ma.Mch
	}
	ma.GetMch()
	return ma.Mch
}


func (ma *MchAccount) GeneratePassword(pwd string) {
	ma.Salt = global.RandStringRunes(10)
	ma.Password = global.PwdPlaintext2CipherText(pwd, ma.Salt)
	return
}

func (ma *MchAccount) LoginSuccess() error {
	now := time.Now()
	ma.Token = sql.NullString{
		String: global.Md5(fmt.Sprintf("%d%s", now.UnixNano(), global.RandStringRunes(12))),
		Valid:  true,
	}
	ma.Status = global.IsYes
	ma.OnlineStatus = global.IsYes
	ma.LastLoginTime = now
	db := config.GetDbDefault()
	if err := db.Save(ma).Error; err != nil {
		return err
	}
	if ma.Mch == nil {
		ma.GetMch()
	}
	ma.Mch.LastLoginTime = now
	db.Save(ma.Mch)
	return nil
}
