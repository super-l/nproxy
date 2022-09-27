package services

type errMsg struct {
	Code int         `json:"code"` // 业务编码
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`          // 错误描述
	ID   string      `json:"id,omitempty"` // 当前请求的唯一ID，便于问题定位，忽略也可以
}

func NewMsg(code int, message string) *errMsg {
	return &errMsg{
		Code: code,
		Msg:  message,
	}
}

var (
	RepeatMsg    = NewMsg(1001, "repeat")
	NoDataMsg    = NewMsg(1004, "no data")
	CommonErrMsg = NewMsg(1005, "error")
)
