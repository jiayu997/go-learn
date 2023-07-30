package kapi

import (
	"gokit/pkg/k8s"
	"gokit/pkg/klog"
	"gokit/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 删除异常pod接口
func DeleteErrorPods(ctx *gin.Context) {
	result, err := k8s.DeleteErrorPods()
	if err != nil {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeError, result+" "+err.Error()))
	} else {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeINFO, "集群异常POD清理成功"))
	}
}

// 查询pod状态接口
func PodStatus(ctx *gin.Context) {
	ws, err := utils.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Abort()
	}
	podlist, err := k8s.GetAllPods()
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, err.Error(), true), time.Now().Add(time.Second))
		ws.Close()
		return
	}
	for _, pod := range podlist.Items {
		var restartCount int
		for _, container := range pod.Status.ContainerStatuses {
			restartCount += int(container.RestartCount)
		}
		var data = gin.H{
			"namespace":  pod.Namespace,
			"podname":    pod.Name,
			"podstatus":  pod.Status.Phase,
			"podrestart": restartCount,
		}
		//ws.WriteMessage(websocket.TextMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, fmt.Sprintf("命名空间：%v POD名称：%v POD状态：%v", pod.Namespace, pod.Name, pod.Status.Phase), false))
		ws.WriteJSON(data)
	}
	ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, "POD状态查询完成", true), time.Now().Add(time.Second))
	ws.Close()
}
