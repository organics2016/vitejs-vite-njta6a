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

func (tty *SSHTty) Connect() {

	//创建ssh
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 4, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            tty.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if len(tty.SecretKey) > 0 {

		signer, err := ssh.ParsePrivateKey(tty.SecretKey)
		if err != nil {
			tty.OutputError(err)
			return
		}
		//config.HostKeyCallback = ssh.FixedHostKey(signer.PublicKey())
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}

	} else if len(tty.Password) > 0 {
		config.Auth = []ssh.AuthMethod{ssh.Password(tty.Password)}
	} else {
		tty.OutputError(errors.New("Not auth"))
		return
	}

	addr := fmt.Sprintf("%s:%d", tty.Host, tty.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		tty.OutputError(err)
		return
	}
	tty.sshClient = sshClient

	session, err := sshClient.NewSession()
	if err != nil {
		tty.OutputError(err)
		return
	}
	tty.sshSession = session

	if stdoutPipe, err := session.StdoutPipe(); err != nil {
		tty.OutputError(err)
		return
	} else {
		tty.readerToWebsocket(stdoutPipe)
	}

	if stderrPipe, err := session.StderrPipe(); err != nil {
		tty.OutputError(err)
		return
	} else {
		tty.readerToWebsocket(stderrPipe)
	}

	if stdinPipe, err := session.StdinPipe(); err != nil {
		tty.OutputError(err)
		return
	} else {
		tty.websocketToWriter(stdinPipe)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		//ssh.ECHO:          0,     // disable echoing
		//ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		//ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		tty.OutputError(err)
		return
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		tty.OutputError(err)
		return
	}

}

func (tty *SSHTty) Close() {
	<-tty.ctx.Done()

	if tty.sshSession != nil {
		tty.sshSession.Close()
	}
	if tty.sshClient != nil {
		tty.sshClient.Close()
	}
	if tty.Websocket.wsConn != nil {
		tty.Websocket.wsConn.Close()
	}
}
