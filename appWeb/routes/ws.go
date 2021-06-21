package routes

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/zqjzqj/instantCustomer/appWeb"
	"github.com/zqjzqj/instantCustomer/appWeb/middleware"
	"github.com/zqjzqj/instantCustomer/logs"
	"github.com/zqjzqj/instantCustomer/ws"
	"net/http"
)

func RegisterWsRoutes(app *iris.Application) {
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		if r.URL.Path == "/wsm" {
			//认证
			ctx := app.ContextPool.Acquire(w, r)
			ma := middleware.RegisterUserAndAuth(ctx)
			if ma == nil {
				app.ContextPool.Release(ctx)
				return
			}
			wsConn, err := ws.NewWsConn(w, r)
			if err != nil {
				ctx.StopWithJSON(200, appWeb.NewResponse(appWeb.ResponseFailCode, "ws connect create fail", nil))
			}
			app.ContextPool.Release(ctx)
			err = ma.ListenWsMsg(wsConn, context.Background())
			if err != nil {
				logs.PrintErr(ma.ID, " listen ws msg err ", err)
			}
			return
		}
		router.ServeHTTP(w, r)
	})
}
