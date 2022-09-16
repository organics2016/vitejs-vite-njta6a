package connection

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"time"
)

type WebSocketSSH struct {
	Request   *http.Request
	Writer    http.ResponseWriter
	Host      string
	Username  string
	Password  string
	SecretKey []byte
	Port      int

	websocket  *websocket.Conn
	sshSession *ssh.Session
	sshClient  *ssh.Client
	active     bool
}

// 升级get请求为webSocket协议
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wss *WebSocketSSH) Connect() error {

	//创建websocket
	ws, err := upGrader.Upgrade(wss.Writer, wss.Request, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	wss.websocket = ws
	wss.active = true

	//创建ssh
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 4, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            wss.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if len(wss.SecretKey) > 0 {

		signer, err := ssh.ParsePrivateKey(wss.SecretKey)
		if err != nil {
			wss.wError(err)
			return errors.WithStack(err)
		}
		//config.HostKeyCallback = ssh.FixedHostKey(signer.PublicKey())
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}

	} else if len(wss.Password) > 0 {
		config.Auth = []ssh.AuthMethod{ssh.Password(wss.Password)}
	} else {
		err := errors.New("Not auth")
		wss.wError(err)
		return errors.WithStack(err)
	}

	addr := fmt.Sprintf("%s:%d", wss.Host, wss.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		wss.wError(err)
		return errors.WithStack(err)
	}
	wss.sshClient = sshClient

	session, err := sshClient.NewSession()
	if err != nil {
		wss.wError(err)
		return errors.WithStack(err)
	}
	wss.sshSession = session

	// websocket主动断开连接
	ws.SetCloseHandler(wss.closeHandler)

	if stdoutPipe, err := session.StdoutPipe(); err != nil {
		wss.wError(err)
		return errors.WithStack(err)
	} else {
		wss.createOutPipe(stdoutPipe)
	}

	if stderrPipe, err := session.StderrPipe(); err != nil {
		wss.wError(err)
		return errors.WithStack(err)
	} else {
		wss.createOutPipe(stderrPipe)
	}

	if stdinPipe, err := session.StdinPipe(); err != nil {
		wss.wError(err)
		return errors.WithStack(err)
	} else {
		wss.createInPipe(stdinPipe)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		//ssh.ECHO:          0,     // disable echoing
		//ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		//ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		wss.wError(err)
		return errors.WithStack(err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		wss.wError(err)
		return errors.WithStack(err)
	}

	// 阻塞，直到websocket断开或容器断开
	session.Wait()

	// 如果websocket先断开连接，这里会重复执行一次，当容器先断开连接时或发生意外，在这里释放资源
	defer wss.closeHandler(0, "")

	return nil
}

func (wss *WebSocketSSH) closeHandler(code int, text string) error {
	wss.active = false
	// 无论任何原因导致的连接关闭，都应该尝试关闭所有已建立的连接
	wss.sshSession.Close()
	wss.sshClient.Close()
	wss.websocket.Close()
	return nil
}

func (wss *WebSocketSSH) createOutPipe(reader io.Reader) {
	go func() {
		for wss.active {
			message := make([]byte, 1*1024)
			_, err := reader.Read(message)
			if err != nil {
				continue
			}
			wss.websocket.WriteMessage(websocket.TextMessage, message)
		}
	}()
}

func (wss *WebSocketSSH) createInPipe(write io.WriteCloser) {
	go func() {
		for wss.active {
			_, message, err := wss.websocket.ReadMessage()
			if err != nil {
				continue
			}
			write.Write(message)
		}
	}()
}

func (wss *WebSocketSSH) wError(err error) {
	wss.websocket.WriteMessage(websocket.TextMessage, []byte(err.Error()))
}
