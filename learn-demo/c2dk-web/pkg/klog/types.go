package klog

// 定义html类型
const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

// 自定义websocekt 返回给前端的响应code
// 定义日志级别
const (
	MessageTypeINFO  = 40001
	MessageTypeDEBUG = 40002
	MessageTypeError = 40003
)

// 日志文件名
const (
	LogFileName = "../ansible.log"
	//LogFileName = "../hosts.ini"
)
