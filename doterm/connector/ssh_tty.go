package connector

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHTty struct {
	Websocket
	Host      string
	Username  string
	Password  string
	SecretKey []byte
	Port      int

	sshSession *ssh.Session
	sshClient  *ssh.Client
}

func (sshTty *SSHTty) Connect() {

	sshTty.initWebSocket(sshTty)
	// 如果websocket先断开连接，这里会重复执行一次，当容器先断开连接时或发生意外，在这里释放资源
	defer sshTty.Close()

	//创建ssh
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 4, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            sshTty.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if len(sshTty.SecretKey) > 0 {

		signer, err := ssh.ParsePrivateKey(sshTty.SecretKey)
		if err != nil {
			sshTty.outputError(err)
			return
		}
		//config.HostKeyCallback = ssh.FixedHostKey(signer.PublicKey())
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}

	} else if len(sshTty.Password) > 0 {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshTty.Password)}
	} else {
		sshTty.outputError(errors.New("Not auth"))
		return
	}

	addr := fmt.Sprintf("%s:%d", sshTty.Host, sshTty.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		sshTty.outputError(err)
		return
	}
	sshTty.sshClient = sshClient

	session, err := sshClient.NewSession()
	if err != nil {
		sshTty.outputError(err)
		return
	}
	sshTty.sshSession = session

	if stdoutPipe, err := session.StdoutPipe(); err != nil {
		sshTty.outputError(err)
		return
	} else {
		sshTty.readerToWebsocket(stdoutPipe)
	}

	if stderrPipe, err := session.StderrPipe(); err != nil {
		sshTty.outputError(err)
		return
	} else {
		sshTty.readerToWebsocket(stderrPipe)
	}

	if stdinPipe, err := session.StdinPipe(); err != nil {
		sshTty.outputError(err)
		return
	} else {
		sshTty.websocketToWriter(stdinPipe)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		//ssh.ECHO:          0,     // disable echoing
		//ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		//ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		sshTty.outputError(err)
		return
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		sshTty.outputError(err)
		return
	}

	// 阻塞，直到websocket断开或容器断开，这种断开视为正常退出
	if err := session.Wait(); err != nil {
	}
}

func (sshTty *SSHTty) Close() {
	sshTty.cancel()

	if sshTty.sshSession != nil {
		sshTty.sshSession.Close()
	}
	if sshTty.sshClient != nil {
		sshTty.sshClient.Close()
	}
	if sshTty.websocket != nil {
		sshTty.websocket.Close()
	}

}