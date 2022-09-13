package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsSSH() *ssh.Client {
	sshHost := "127.0.0.1"
	sshUsername := "vagrant"
	sshPassword := "xxxxxx"
	sshType := "key"                                                            //password 或者 key
	sshKeyPath := "D:/vagrant/.vagrant/machines/default/virtualbox/private_key" //ssh id_rsa.id 路径"
	sshPort := 2222

	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            sshUsername,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		// HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}

	return sshClient
}

type websocketSSH struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *websocketSSH) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	key, err := os.ReadFile(kPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func ping(c *gin.Context) {

	sshClient := wsSSH()

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	defer session.Close()

	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			message := make([]byte, 1*1024)
			_, err := stdoutPipe.Read(message)
			if err != nil {
				continue
			}
			fmt.Printf("send : %s\n", string(message))
			ws.WriteMessage(websocket.TextMessage, message)
		}
	}()

	stderrPipe, err := session.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			message := make([]byte, 4*1024)
			_, err := stderrPipe.Read(message)
			if err != nil {
				continue
			}
			fmt.Printf("err : %s\n", string(message))
			ws.WriteMessage(websocket.TextMessage, message)
		}
	}()

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				continue
			}
			fmt.Printf("receive : %s\n", string(message))
			stdinPipe.Write(message)
		}
	}()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		//ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	session.Wait()
}

func main() {
	// 125800000
	// 0.0340246422893482
	// 8,690
	// 0.004352146263911
	r := gin.Default()
	r.GET("/ping", ping)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
