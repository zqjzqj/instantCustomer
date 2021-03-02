package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/zqjzqj/instantCustomer/appWeb/controller"
)

func RegisterApiRoutes(app *iris.Application) {
	mvc.Configure(app.Party("/"), func(application *mvc.Application) {
		application.Handle(&controller.SiteController{})
	})

}
