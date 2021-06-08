package global

import "github.com/kataras/iris/v12"

const (
	ReqTokenName       = "token"
	ReqTokenHeaderName = "X-Token"
)

func GetReqToken(ctx iris.Context) string {
	token := ctx.GetHeader(ReqTokenHeaderName)
	if token == "" {
		token = ctx.URLParamTrim(ReqTokenName)
	}
	if token == "" {
		token = ctx.PostValue(ReqTokenName)
	}
	return token
}
