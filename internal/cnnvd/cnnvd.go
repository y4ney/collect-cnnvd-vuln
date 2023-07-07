package cnnvd

const (
	Schema = "https"
	Domain = "www.cnnvd.org.cn"
)

// ResCode 响应码
type ResCode struct {
	Code    int    `json:"code,omitempty"`    // 代码
	Success bool   `json:"success,omitempty"` // 是否成功
	Message string `json:"message,omitempty"` // 消息
	Time    string `json:"time,omitempty"`    // 时间
}
