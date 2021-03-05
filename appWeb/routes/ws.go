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
		ws := websocket.New(websocket.DefaultGobwasUpgrader, application)/*

		//修改连接ID
		ws.IDGenerator = func(w http.ResponseWriter, r *http.Request) string {
			token := r.URL.Query().Get(global.ReqTokenName)
			if token == "" {
				token = r.Header.Get(global.ReqTokenHeaderName)
			}
			ma := &models.MchAccount{}
			config.GetDbDefault().Where("token = ?", token).Find(ma)
			return strconv.FormatUint(ma.ID, 10)
		}*/
		application.Router.Get("/", websocket.Handler(ws))
	})
}