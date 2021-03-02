package appWeb

const (
	ResponseSuccessCode = 0
	ResponseFailCode = 1
	ResponseNotLoginCode = -1
)

type ResponseFormat struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse(code int, msg string, data interface{}) *ResponseFormat {
	if msg == "" {
		if code == ResponseSuccessCode {
			msg = "操作成功"
		} else if code == ResponseFailCode {
			msg = "操作失败"
		} else if code == ResponseNotLoginCode {
			msg = "账户未登录"
		}
	}
	return &ResponseFormat{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
