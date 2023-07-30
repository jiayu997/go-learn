package klog

import "github.com/gorilla/websocket"

// 用于设置websocket 发送消息时，传递给前端的响应code
func GenerateHtmlContent(code int, content string, self bool) []byte {
	switch code {
	case MessageTypeDEBUG:
		if self {
			return websocket.FormatCloseMessage(MessageTypeDEBUG, content)
		} else {
			return []byte(content)
		}
	case MessageTypeINFO:
		if self {
			return websocket.FormatCloseMessage(MessageTypeINFO, content)
		} else {
			return []byte(content)
		}
	case MessageTypeError:
		if self {
			return websocket.FormatCloseMessage(MessageTypeError, content)
		} else {
			return []byte(content)
		}
	}
	return []byte{}
}

func GenerateHtml(code int, content string) []byte {
	switch code {
	case MessageTypeDEBUG:
		return []byte("<p id='debug'>" + "[DEBUG] " + content + "</p>")
	case MessageTypeINFO:
		return []byte("<p id='info'>" + "[INFO] " + content + "</p>")
	case MessageTypeError:
		return []byte("<p id='error'>" + "[ERROR] " + content + "</p>")
	}
	return []byte{}
}
