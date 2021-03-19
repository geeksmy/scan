package blasting

import (
	"time"

	"golang.org/x/crypto/ssh"
)

func NewConnSSH(addr, user, pass string, timeout int) bool {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         time.Duration(timeout) * time.Second,
	}
	// 建立与SSH服务器的连接
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return false
	}
	defer sshClient.Close()
	return true
}
