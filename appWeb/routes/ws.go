package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/appWeb/appRegister"
	"github.com/zqjzqj/instantCustomer/appWeb/middleware"
	"github.com/zqjzqj/instantCustomer/appWeb/wsController"
	"github.com/zqjzqj/instantCustomer/services"
)


func RegisterWebsocketRoutes(app *iris.Application) {
	app.Logger().SetLevel("debug")

	//商户的客服ws空间
	mvc.Configure(app.Party("/ws/mch",  middleware.MchAccountAuth), func(application *mvc.Application) {
		application.Register(appRegister.MchAccount)
		application.HandleWebsocket(&wsController.MchWsController{Namespace:services.WsNamespaceMch})
		ws := websocket.New(websocket.DefaultGobwasUpgrader, application)
		application.Router.Get("/", websocket.Handler(ws))
	})
}