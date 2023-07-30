package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"gokit/pkg/ssh"

	"github.com/gorilla/websocket"
	"github.com/spf13/pflag"
)

type MonitorInfo struct {
	IP     string `json:"monitor"`
	Enable bool   `json:"enable"`
}

type LogInfo struct {
	IP     []string `json:"log"`
	Enable bool     `json:"enable"`
}

type VipInfo struct {
	IP     string `json:"vip"`
	Enable bool   `json:"enable"`
}

type BackupInfo struct {
	IP     string `json:"backup"`
	Enable bool   `json:"enable"`
}

type amp struct {
	Version            string   `json:"version" binding:"required"`
	Master             string   `json:"master" binding:"required"`
	MasterControlPlane []string `json:"masterControl"`
	Node               []string
	K8SVIP             string   `json:"k8svip"`
	NfsHarbor          string   `json:"nfsharbor"`
	ClusterAdmin       []string `form:"clusteradmin"`
	Monitor            MonitorInfo
	Log                LogInfo
	BusinVIP           VipInfo
	Backup             BackupInfo
	NetworkInterface   string
	SshPort            string
	SshPassword        string
}

var (
	Amp      = new(amp)
	UpGrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	JwtToken string // jwt token
	Auth     bool
)

func (amp *amp) InitFlag() error {
	pflag.StringVar(&amp.SshPort, "port", "22", "--port=22")
	pflag.StringVar(&amp.SshPassword, "password", "", "--password=password")
	pflag.Parse()
	if amp.SshPassword == "" {
		return errors.New(pflag.CommandLine.FlagUsages())
	}

	// 本地SSH PASSWORD检查
	err := ssh.Ping("localhost", amp.SshPassword, amp.SshPort)
	if err != nil {
		return err
	}
	return nil
}

// 结构体序列化，便于查看信息
func (amp *amp) Print() {
	//js, err := json.Marshal(amp)
	js, err := json.MarshalIndent(amp, "", "\t")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s", js)
}

// 获取 master节点，本地网卡名称
// 本地网卡名称检测
func (amp *amp) InitNetworkInterface() error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// 获取所有的网卡信息
	for _, inter := range interfaces {
		addrlist, err := inter.Addrs()
		if err != nil {
			return err
		}
		// 获取网卡上的所有IP(IPV4/IPV6)
		for _, ip := range addrlist {
			if strings.Split(ip.String(), "/")[0] == amp.Master {
				amp.NetworkInterface = inter.Name
				return nil
			}
		}
	}
	return errors.New("master ip not found")
}
