package appRegister

import (
	"github.com/kataras/iris/v12"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/models"
)

//该注册应用必须在中间件MchAccountAuth之后 否则将会获取不到mAccount
func MchAccount(ctx iris.Context) *models.MchAccount {
	var mch *models.MchAccount
	mch, ok := ctx.Values().Get("mAccount").(*models.MchAccount)
	if ok {
		return mch
	}
	token := global.GetReqToken(ctx)
	if token == "" {
		return nil
	}
	config.GetDbDefault().Where("token = ?", token).Find(mch)
	return mch
}
