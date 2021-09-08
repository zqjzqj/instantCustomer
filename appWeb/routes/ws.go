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
			HandleWsm(app, w, r)
			return
		}
		if r.URL.Path == "/wsc" {
			HandleWsm(app, w, r)
			return
		}
		if r.URL.Path == "/wsv" {
			HandleWsVisitor(app, w, r)
			return
		}
		router.ServeHTTP(w, r)
	})
}

func HandleWsVisitor(app *iris.Application, w http.ResponseWriter, r *http.Request) {
	ctx := app.ContextPool.Acquire(w, r)
	wsConn, err := ws.NewWsConn(w, r)
	if err != nil {
		logs.PrintErr("ws connect fail ", err)
		app.ContextPool.Release(ctx)
		ctx.StopWithJSON(200, appWeb.NewResponse(appWeb.ResponseFailCode, "ws connect create fail", nil))
		return
	}
	app.ContextPool.Release(ctx)
	writeChan := make(chan *ws.Message)
	go func() {
		for msg := range writeChan {
			err := wsConn.WriteMsg(msg)
			if err != nil {
				logs.PrintErr("write err ", err)
			}
		}
	}()
	err = wsConn.ListenMsg(writeChan)
	if err != nil {
		logs.PrintErr(wsConn.ID, " listen ws msg err ", err)
	}
}

func HandleWsm(app *iris.Application, w http.ResponseWriter, r *http.Request) {
	//认证
	ctx := app.ContextPool.Acquire(w, r)
	ma := middleware.RegisterUserAndAuth(ctx)
	if ma == nil {
		logs.PrintErr("user auth fail")
		app.ContextPool.Release(ctx)
		return
	}
	wsConn, err := ws.NewWsConn(w, r)
	if err != nil {
		logs.PrintErr("ws connect fail ", err)
		app.ContextPool.Release(ctx)
		ctx.StopWithJSON(200, appWeb.NewResponse(appWeb.ResponseFailCode, "ws connect create fail", nil))
		return
	}
	app.ContextPool.Release(ctx)
	err = ma.ListenWsMsg(wsConn, context.Background())
	if err != nil {
		logs.PrintErr(ma.ID, " listen ws msg err ", err)
	}
}

func HandleWsc(app *iris.Application, w http.ResponseWriter, r *http.Request) {

}
