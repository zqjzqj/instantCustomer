package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/zqjzqj/instantCustomer/appWeb"
	"github.com/zqjzqj/instantCustomer/appWeb/appRegister"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/models"
)

func MchAccountAuth(ctx iris.Context) {
	token := global.GetReqToken(ctx)
	if token == "" {
		ctx.StopWithJSON(200, appWeb.NewResponse(appWeb.ResponseNotLoginCode, "", nil))
		return
	}

	ma := &models.MchAccount{}
	config.GetDbDefault().Where("token = ?", token).Find(ma)
	if ma.ID == 0 {
		ctx.StopWithJSON(200, appWeb.NewResponse(appWeb.ResponseNotLoginCode, "", nil))
		return
	}
	ctx.Values().Set(appRegister.MAccountKey, ma)
	ctx.Next()
}