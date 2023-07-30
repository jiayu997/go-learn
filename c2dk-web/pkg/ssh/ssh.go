package ssh

import (
	"time"

	"golang.org/x/crypto/ssh"
)

// ping 网络测试
func Ping(host, password, port string) error {
	// ssh 配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	// 获取ssh client
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return err
	}
	client.Close()
	return nil
}

// 获取SSH Client
func GetSSHClient(host, password, port string) (*ssh.Client, error) {
	// ssh 配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	// 获取ssh client
	sshclient, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, err
	}
	return sshclient, nil
}

// 远程执行命令
func RunCmd(client *ssh.Client, cmd string) ([]byte, error) {
	session, err := client.NewSession()
	defer client.Close()
	defer session.Close()
	if err != nil {
		return []byte(err.Error()), err
	}
	by, err := session.CombinedOutput(cmd)
	if err != nil {
		return []byte(err.Error()), err
	}
	return by, nil
}
