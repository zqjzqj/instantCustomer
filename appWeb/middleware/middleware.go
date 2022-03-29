package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/zqjzqj/instantCustomer/appWeb"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/models"
	"net/http"
)

func RegisterUserAndAuth(ctx iris.Context) *models.MchAccount {
	token := global.GetReqToken(ctx)
	if token == "" {
		ctx.StopWithJSON(http.StatusOK, appWeb.NewResponse(appWeb.ResponseNotLoginCode, "", nil))
		return nil
	}

	ma := &models.MchAccount{}
	global.GetDb().Where("token = ?", token).Find(ma)
	if ma.ID == 0 {
		ctx.StopWithJSON(http.StatusOK, appWeb.NewResponse(appWeb.ResponseNotLoginCode, "", nil))
		return nil
	}
	return ma
}
