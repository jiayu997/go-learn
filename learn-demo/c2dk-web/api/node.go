package kapi

import (
	"bufio"
	"fmt"
	"gokit/pkg/klog"
	"gokit/utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"gokit/pkg/k8s"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func DeleteNode(ctx *gin.Context) {
	ip, exist := ctx.GetQuery("addr")
	if !exist {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeError, "未携带参数"))
		return
	}

	// 判断当前集群是否正常（防止集群未搭建就删除节点）
	_, err := k8s.GetClient()
	if err != nil {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeError, err.Error()))
		return
	}

	err = k8s.DeleteNode(ip)
	if err != nil {
		ctx.Writer.Header().Set("Content-Type", "text/html")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeError, err.Error()))
		return
	}
	ctx.Writer.Header().Set("Content-Type", "text/html")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(klog.GenerateHtml(klog.MessageTypeINFO, ip+"节点删除成功"))
}

func AddNode(ctx *gin.Context) {
	ws, err := utils.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Abort()
	}
	defer ws.Close()
	// 判断是否存在hosts.ini文件
	if !utils.FileExist("../hosts.ini") {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "hosts.ini文件不存在", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	// 获取websocket 发送过来的数据
	_, data, err := ws.ReadMessage()
	ip := string(data)
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "提交的数据异常", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	// 判断当前集群是否正常（防止集群未搭建就添加节点）
	_, err = k8s.GetClient()
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "K8S集群连接失败", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	// 判断集群是否已存在该节点
	_, err = k8s.SearchNode(ip)
	if err == nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, ip+"节点已存在，请勿重复添加节点", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	cmd1 := exec.Command("/bin/bash", "-c", `sed -ri '/\[newnode\]/,$d' ../hosts.ini`)
	_, err = cmd1.CombinedOutput()
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "配置生成失败", true), time.Now().Add(time.Second))
		fmt.Println(err)
		ctx.Abort()
		return
	}

	host_f, err := os.OpenFile("../hosts.ini", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "配置生成失败", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	// 写入newnode配置
	defer host_f.Close()
	_, err = host_f.WriteString("[newnode]\n" + ip + "\n")
	if err != nil {
		fmt.Println(err)
	}

	//执行Ansible
	//	cmd := exec.Command("/bin/bash", "-c", "../c2dkctl c2dkweb newnode "+ip +" " +)
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("../c2dkctl c2dkweb newnode %s %s %s", ip, utils.Amp.SshPassword, utils.Amp.SshPort))
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	readout := bufio.NewReader(stdout)
	readerr := bufio.NewReader(stderr)
	go func() {
		for {
			line, err := readout.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				//ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
				break
			}
			if err == nil {
				ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeINFO, line))
			}
		}
	}()
	go func() {
		for {
			line, err := readerr.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				break
				//ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
			}
			if err == nil {
				ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, line))
			}
		}
	}()
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "自动化执行出错，程序退出码为: "+fmt.Sprintf("%d", exiterr.ExitCode()), true), time.Now().Add(time.Second))
			ws.Close()
			return
		}
	}
	ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, ip+"节点添加成功", true), time.Now().Add(time.Second))
}
