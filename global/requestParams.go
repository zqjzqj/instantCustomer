package global

import "github.com/kataras/iris/v12"

func GetReqToken(ctx iris.Context) string {
	token := ctx.GetHeader("X-token")
	if token == "" {
		token = ctx.URLParamTrim("token")
	}
	if token == "" {
		token = ctx.PostValue("token")
	}
	return token
}
