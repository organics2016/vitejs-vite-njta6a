package connector

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHTty struct {
	*Websocket
	Host     string
	Port     int
	Username string
	Password string

	PriKey []byte

	sshSession *ssh.Session
	sshClient  *ssh.Client
}

func (tty *SSHTty) Connect() error {

	config := &ssh.ClientConfig{
		Timeout:         time.Second * 4,
		User:            tty.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if len(tty.PriKey) > 0 {

		privateKey, err := ssh.ParsePrivateKey(tty.PriKey)
		if err != nil {
			return err
		}

		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(privateKey)}

	} else if len(tty.Password) > 0 {
		config.Auth = []ssh.AuthMethod{ssh.Password(tty.Password)}
	} else {
		return errors.New("Not auth")
	}

	addr := fmt.Sprintf("%s:%d", tty.Host, tty.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	tty.sshClient = sshClient

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	tty.sshSession = session

	if stdoutPipe, err := session.StdoutPipe(); err != nil {
		return err
	} else {
		tty.readerToWebsocket(stdoutPipe)
	}

	if stderrPipe, err := session.StderrPipe(); err != nil {
		return err
	} else {
		tty.readerToWebsocket(stderrPipe)
	}

	if stdinPipe, err := session.StdinPipe(); err != nil {
		return err
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
		return err
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		return err
	}

	<-tty.ctx.Done()

	return nil
}

func (tty *SSHTty) Close() {

	if tty.sshSession != nil {
		tty.sshSession.Close()
	}
	if tty.sshClient != nil {
		tty.sshClient.Close()
	}
}
