package kapi

import (
	"bufio"
	"fmt"
	"gokit/pkg/k8s"
	"gokit/pkg/klog"
	"gokit/pkg/ssh"
	"gokit/utils"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Deploy(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "deploy.html", nil)
}

// interface：/api/v1/deploycheck
// method: GET
func DeployCheck(ctx *gin.Context) {
	ws, err := utils.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Abort()
	}
	defer ws.Close()

	// 序列化websocket发送过来的数据
	err = ws.ReadJSON(&utils.Amp)
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "提交检查数据序列化失败", true), time.Now().Add(time.Second))
		return
	}
	utils.Amp.Print()

	// master节点根据IP反查网卡,并赋值
	err = utils.Amp.InitNetworkInterface()
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "主节点网卡/网络检测异常", true), time.Now().Add(time.Second))
		return
	} else {
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, "主节点状态检测通过", false))
	}

	// master ssh 网络连通性与网卡检测
	err = ssh.Ping(utils.Amp.Master, utils.Amp.SshPassword, utils.Amp.SshPort)
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, utils.Amp.Master+"节点SSH连接失败", true), time.Now().Add(time.Second))
		return
	} else {
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, utils.Amp.Master+"节点网络连通性与网卡检测通过", false))
	}

	// master-control ssh校验与网卡校验
	for _, node := range utils.Amp.MasterControlPlane {
		client, err := ssh.GetSSHClient(node, utils.Amp.SshPassword, utils.Amp.SshPort)
		if err != nil {
			ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, node+"节点SSH连接失败", true), time.Now().Add(time.Second))
			return
		}
		_, err = ssh.RunCmd(client, "ip link show | grep "+utils.Amp.NetworkInterface)
		if err != nil {
			ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, node+"节点网卡与主节点网卡名不一致", true), time.Now().Add(time.Second))
			return
		}
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, node+"节点网络连通性与网卡检测通过", false))
		client.Close()
	}

	// node ssh校验与网卡校验
	for _, node := range utils.Amp.Node {
		client, err := ssh.GetSSHClient(node, utils.Amp.SshPassword, utils.Amp.SshPort)
		if err != nil {
			ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, node+" 节点网卡与主节点网卡名不一致", true), time.Now().Add(time.Second))
			ws.Close()
			return
		}
		_, err = ssh.RunCmd(client, "ip link show | grep "+utils.Amp.NetworkInterface)
		if err != nil {
			ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, node+" 节点网卡与主节点网卡名不一致", true), time.Now().Add(time.Second))
			ws.Close()
			return
		}
		client.Close()
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, node+"节点网络连通性与网卡检测通过", false))
	}

	// 生成模板解析
	err = ampGenerateConfig(ws)
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "Ansible配置文件生成失败！", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, "节点检查正常，具备部署条件", true), time.Now().Add(time.Second))
	ws.Close()
}

// interface：/api/v1/deployrun
// method: GET
func DeployRun(ctx *gin.Context) {
	ws, err := utils.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Abort()
		return
	}
	defer ws.Close()

	// 执行Ansible
	//cmd := exec.Command("/bin/bash", "-c", "../c2dkctl c2dkweb generateConfig "+utils.Amp.NetworkInterface+" "+utils.Amp.SshPassword+" "+utils.Amp.SshPort+" && ../c2dkctl c2dkweb install")
	cmd := exec.Command("/bin/bash", "-c", "../c2dkctl c2dkweb install "+utils.Amp.NetworkInterface+" "+utils.Amp.SshPassword+" "+utils.Amp.SshPort)
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

	ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeINFO, "amp2.9 平台管控台访问地址为："+"http://"+strings.Split(ctx.Request.Host, ":")[0]+":30000"))
	ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, "自动化执行成功！可以访问平台了", true), time.Now().Add(time.Second))
}

// interface：/api/v1/deploylog
// method: GET
func DeployLog(ctx *gin.Context) {
	exist := utils.FileExist(klog.LogFileName)
	if exist {
		ctx.Header("Content-Disposition", "attachment; filename="+time.Now().Format("2006-01-02")+"-C2DK部署日志")
		ctx.Header("Content-Type", "application/text/plain")
		ctx.File(klog.LogFileName)
	} else {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.WriteString("当前未生成日志文件，下载失败")
		ctx.Writer.Flush()
	}
}

// interface: /api/v1/deploycomponent
// method: GET
func DeployComponent(ctx *gin.Context) {
	ws, err := utils.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Abort()
		return
	}
	defer ws.Close()

	// K8S连接检查
	_, err = k8s.GetClient()
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "K8S集群异常", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	// 获取websocket 发送过来的数据
	_, data, err := ws.ReadMessage()
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "获取组件异常", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	component := string(data)
	re, err := utils.CompareString(component, "c2cloud|monitor|harbor")
	if err != nil {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, "获取组件异常", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}
	if !re {
		ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeError, component+"该组件不支持", true), time.Now().Add(time.Second))
		ctx.Abort()
		return
	}

	// 执行Ansible
	cmd := exec.Command("/bin/bash", "-c", "../c2dkctl tags "+component)
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
	ws.WriteControl(websocket.CloseMessage, klog.GenerateHtmlContent(klog.MessageTypeINFO, component+"组件部署成功", true), time.Now().Add(time.Second))
}

func ampGenerateConfig(ws *websocket.Conn) error {
	// 渲染配置文件
	tmpl, err := template.New("").Delims("[[", "]]").ParseFiles("./static/template/hosts.ini.tmpl", "./static/template/all.yml.tmpl", "./static/template/amp2.9.yaml.tmpl")
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
		return err
	}

	host_f, err := os.OpenFile("../hosts.ini", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
		return err
	}

	all_f, err := os.OpenFile("../group_vars/all.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
		return err
	}

	if utils.Amp.Version == "2.9" {
		amp_f, err := os.OpenFile("../group_vars/amp2.9.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
			return err
		}
		defer amp_f.Close()
		err = tmpl.ExecuteTemplate(amp_f, "amp29", utils.Amp)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
			return err
		}
	} else if utils.Amp.Version == "3.0" {
		amp_f, err := os.OpenFile("../group_vars/amp3.0.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
			return err
		}
		defer amp_f.Close()
		err = tmpl.ExecuteTemplate(amp_f, "amp30", utils.Amp)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
			return err
		}
	} else if utils.Amp.Version == "lite" {
		amp_f, err := os.OpenFile("../group_vars/amplite.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
			return err
		}
		defer amp_f.Close()
		err = tmpl.ExecuteTemplate(amp_f, "amplite", utils.Amp)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
			return err
		}
	}
	defer host_f.Close()
	defer all_f.Close()

	err = tmpl.ExecuteTemplate(host_f, "host", utils.Amp)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
		return err
	}
	err = tmpl.ExecuteTemplate(all_f, "all", utils.Amp)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, klog.GenerateHtml(klog.MessageTypeError, err.Error()))
		return err
	}
	return nil
}
