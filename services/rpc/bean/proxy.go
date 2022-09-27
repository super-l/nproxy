package bean

type ProxyReply struct {
	Code int         `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

func (pr *ProxyReply) Error(code int, msg string) {
	pr.Code = code
	pr.Msg = msg
}

func (pr *ProxyReply) Success(code int, data interface{}, msg string) {
	pr.Code = code
	pr.Data = data
	pr.Msg = msg
}

type GetProxyArgs struct {
	ProtocolType string `json:"protocol_type,omitempty"` // 协议类型
	LineType     int    `json:"line_type,omitempty"`     // 线路类型
	Country      string `json:"country,omitempty"`       // 国家
	Count        int    `json:"count,omitempty"`         // 个数
}

type AddProxyArgs struct {
	ProtocolType string `json:"protocol_type"` // 协议类型
	LineType     int    `json:"line_type"`     // 线路类型
	Value        string `json:"value"`         // 值
	Source       string `json:"source"`        // 来源
}

type UpdateProxyArgs struct {
	Id           int    `json:"id,omitempty"`
	ProtocolType string `json:"protocol_type,omitempty"` // 协议类型
	LineType     int    `json:"line_type,omitempty"`     // 线路类型
	Value        string `json:"value,omitempty"`         // 值
	Source       string `json:"source"`                  // 来源
}

type ListProxyArgs struct{}

type DeleteProxyArgs struct {
	IdList []string
}

type DeleteProxyApiArgs = DeleteProxyArgs

type AddProxyApiArgs struct {
	ProtocolType string `json:"protocol_type"` // 协议类型
	LineType     int    `json:"line_type"`     // 线路类型
	Value        string `json:"value"`         // 值
}

type UpdateProxyApiArgs struct {
	Id           int    `json:"id"`
	ProtocolType string `json:"protocol_type"` // 协议类型
	LineType     int    `json:"line_type"`     // 线路类型
	Value        string `json:"value"`         // 值
}

type ListProxyApiArgs struct{}
