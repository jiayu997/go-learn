package kapi

import (
	"gokit/pkg/k8s"
	"gokit/pkg/klog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DnsDelay(ctx *gin.Context) {
	server, err := k8s.GetDNSER("kube-system", "kube-dns")
	if err != nil {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeError, server+" "+err.Error()))
	}

	result := k8s.DnsCheck("kube-dns.kube-system.svc.cluster.local", server)
	ctx.Writer.Header().Set("Content-Type", "text/html")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeINFO, result))
}
