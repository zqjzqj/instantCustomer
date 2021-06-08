package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zqjzqj/instantCustomer/appWeb"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/models"
)

type SiteController struct {
}

func (s *SiteController) BeforeActivation(b mvc.BeforeActivation) {}

//登录
func (s *SiteController) PostLogin(ctx iris.Context) *appWeb.ResponseFormat {
	ma, err := models.LoginMch(ctx.PostValue("phone"), ctx.PostValue("password"))
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	ma.GetMch()
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "登录成功", map[string]interface{}{
		"mch_id":          ma.MchId,
		"token":           ma.Token.String,
		"phone":           ma.Phone,
		"real_name":       ma.RealName,
		"mch_name":        ma.Mch.Name,
		"last_login_time": ma.LastLoginTime.Format(global.DateTimeFormatStr),
	})
}

//注册
func (s *SiteController) PostRegister(ctx iris.Context) *appWeb.ResponseFormat {
	ma, err := models.NewCreateMerchant(ctx.PostValue("phone"), ctx.PostValue("name"), ctx.PostValue("password"))
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	//自动登录
	_ = ma.LoginSuccess()
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "注册成功", map[string]interface{}{
		"mch_id":          ma.MchId,
		"token":           ma.Token.String,
		"phone":           ma.Phone,
		"real_name":       ma.RealName,
		"mch_name":        ma.Mch.Name,
		"last_login_time": ma.LastLoginTime.Format(global.DateTimeFormatStr),
	})
}
