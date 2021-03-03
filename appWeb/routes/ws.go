package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/appWeb/appRegister"
	"github.com/zqjzqj/instantCustomer/appWeb/controller"
	"github.com/zqjzqj/instantCustomer/appWeb/middleware"
)

func RegisterWebsocketRoutes(app *iris.Application) {
	app.Logger().SetLevel("debug")
	mvc.Configure(app.Party("/websocket",  middleware.MchAccountAuth), func(application *mvc.Application) {
		application.Register(appRegister.MchAccount)
		application.HandleWebsocket(&controller.WebsocketController{Namespace:"de"})
		ws := websocket.New(websocket.DefaultGobwasUpgrader, application)
		application.Router.Get("/", websocket.Handler(ws))
	})
}