package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zqjzqj/instantCustomer/appWeb/appRegister"
	"github.com/zqjzqj/instantCustomer/appWeb/controller"
	"github.com/zqjzqj/instantCustomer/appWeb/middleware"
)

func RegisterApiRoutes(app *iris.Application) {
	mvc.Configure(app.Party("/"), func(application *mvc.Application) {
		application.Handle(&controller.SiteController{})
	})

	//这里需要注册检查商户登录的中间件
	mvc.Configure(app.Party("/merchant", middleware.MchAccountAuth), func(application *mvc.Application) {
		application.Register(appRegister.MchAccount)
		application.Handle(&controller.MerchantController{})
	})
}
