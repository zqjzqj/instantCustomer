package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/logs"
	"github.com/zqjzqj/instantCustomer/sErr"
	"github.com/zqjzqj/instantCustomer/ws"
	"gorm.io/gorm"
	"time"
)

const (
	MchAccountMaxSessionDefault = 99

	//角色常量
	MchAccountRoleAdmin    = 2
	MchAccountRoleCustomer = 1
)

type MchAccount struct {
	FieldsModel
	MchId         uint64         `gorm:"index:idx_mchId;not null;default:0"`
	Phone         string         `gorm:"size:20;comment:手机号;not null;index:idx_phone,unique"`
	Password      string         `gorm:"type:char(32);default:'';comment:密码md5"`
	Salt          string         `gorm:"type:varchar(32);default:'';comment:盐"`
	Avatar        string         `gorm:"type:text;comment:头像"`
	WxOpenId      sql.NullString `gorm:"size:200;comment:微信openid;index:idx_wx,unique"`
	RealName      string         `gorm:"size:20;comment:姓名;default:''"`
	MaxSession    uint           `gorm:"type:int(11) unsigned;default:99;comment:最大会话数"`
	Role          uint8          `gorm:"default:1;comment:1普通客服 2管理员"`
	Token         sql.NullString `gorm:"size:32;index:idx_token,unique';comment:用户移动端登陆token"`
	TokenStatus   uint8          `gorm:"default:1;comment:0禁用 1正常"`
	Status        uint8          `gorm:"default:1;comment:0禁用 1正常"`
	OnlineStatus  uint8          `gorm:"default:0;comment:商户在线状态 1在线 0离线 2隐身"`
	LastLoginTime time.Time      `gorm:"comment:最近一次登陆时间"`
	FieldsExtendsJsonType
	ConnId     sql.NullString `gorm:"size:200;index:idx_connId,unique;comment:ws的连接id"` //当前连接id
	Mch        *Merchant      `gorm:"foreignKey:Id;-"`
	SessionNum int            `gorm:"default:0;comment:当前会话数"`
}

func (ma *MchAccount) TableName() string {
	return "mch_account"
}

func LoginMch(phone, pwd string) (*MchAccount, error) {
	ma := &MchAccount{}
	global.GetDb().Where("phone = ?", phone).Find(ma)
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
		MchId:      mch.ID,
		Phone:      phone,
		Password:   "",
		Salt:       "",
		Avatar:     "",
		RealName:   realName,
		MaxSession: MchAccountMaxSessionDefault,
		Role:       role,
		Token: sql.NullString{
			String: "",
			Valid:  false,
		},
		TokenStatus:  global.IsNo,
		Status:       global.IsYes,
		OnlineStatus: global.IsNo,
		WxOpenId: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	ma.GeneratePassword(password)
	if tx == nil {
		tx = global.GetDb()
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
	global.GetDb().Where("id = ?", ma.MchId).Find(ma.Mch)
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
	ma.TokenStatus = global.IsYes
	ma.Status = global.IsYes
	ma.OnlineStatus = global.IsYes
	ma.LastLoginTime = now
	db := global.GetDb()
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

func (ma *MchAccount) IsOnline() bool {
	_, ok := ws.FindWsConn(ma.ConnId.String)
	if ok {
		if ma.OnlineStatus == global.OnlineStatusLeave {
			ma.OnlineStatus = global.OnlineStatusYes
			global.GetDb().Table(ma.TableName()).Where("id = ?", ma.ID).Updates(map[string]interface{}{
				"online_status": ma.OnlineStatus,
			})
		}
		return true
	}
	if ma.OnlineStatus == global.OnlineStatusYes || ma.OnlineStatus == global.OnlineStatusHide {
		ma.OnlineStatus = global.OnlineStatusHide
		global.GetDb().Table(ma.TableName()).Where("id = ?", ma.ID).Updates(map[string]interface{}{
			"conn_id":       nil,
			"online_status": ma.OnlineStatus,
		})
	}
	return false
}

func (ma *MchAccount) ListenWsMsg(wsConn *ws.WsConn, ctx context.Context) error {
	if ma.IsOnline() {
		err := ws.CloseById(ma.ConnId.String)
		logs.PrintlnWarning("close connect ", ma.ConnId.String, ma.ID)
		if err != nil {
			return err
		}
	}
	ma.ConnId = sql.NullString{
		String: wsConn.ID(),
		Valid:  true,
	}
	ma.OnlineStatus = global.OnlineStatusYes
	//这里更新一下在线状态
	if global.GetDb().Table(ma.TableName()).Where("id = ?", ma.ID).Updates(map[string]interface{}{
		"conn_id":       ma.ConnId,
		"online_status": ma.OnlineStatus,
	}).RowsAffected == 0 {
		return sErr.New("update conn id fail")
	}
	msgChan := make(chan *ws.Message)
	c, cc := context.WithCancel(ctx)
	defer func() {
		ma.ConnId = sql.NullString{
			String: "",
			Valid:  false,
		}
		ma.OnlineStatus = global.OnlineStatusLeave
		global.GetDb().Table(ma.TableName()).Where("id = ?", ma.ID).Updates(map[string]interface{}{
			"conn_id":       nil,
			"online_status": global.OnlineStatusLeave,
		})
		cc()
	}()
	go ws.HandleMsgForwardToClient(msgChan, c, nil)

	return wsConn.ListenMsg(msgChan)
}
