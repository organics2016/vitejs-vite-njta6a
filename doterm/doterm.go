package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"organics.ink/doterm/cloud"
	"organics.ink/doterm/connector"
	"os"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
	Error() string
}

func readKey(filePath string) []byte {
	key, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return key
}

type Connection struct {
	Env   string `form:"env" binding:"required"`
	Token string `form:"token" binding:"required"`
}

func doterm(c *gin.Context) {

	websocket := connector.InitWebSocket(c.Writer, c.Request, context.Background())

	var connection Connection
	if err := c.ShouldBindQuery(&connection); err != nil {
		websocket.OutputError(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if connection.Env == "local" {
		// TODO
		return
	}

	auth := cloud.JianmuAuthorize{}
	connData := auth.Authorize()

	var tty connector.Manager

	switch connData.Type {
	case "host":
		tty = &connector.SSHTty{
			Websocket: websocket,

			Host:      "127.0.0.1",
			Username:  "vagrant",
			SecretKey: readKey("D:/vagrant/.vagrant/machines/default/virtualbox/private_key"),
			Port:      2222,
		}
		break
	case "docker":
		tty = &connector.DockerTty{
			Websocket: websocket,

			Host:        "tcp://127.0.0.1:2375",
			ContainerID: "62c41d9cf865b22ba5de8e45462b5744ae34ffd056dbab48542ff1e48c690678",
		}
		break
	case "kubernetes":
		tty = &connector.K8STty{
			Websocket: websocket,

			Host:         "https://127.0.0.1:49154",
			PodNamespace: "default",
			PodName:      "shell-demo",
			CertData:     readKey("D:/vagrant/.minikube/profiles/multinode-demo/client.crt"),
			KeyData:      readKey("D:/vagrant/.minikube/profiles/multinode-demo/client.key"),
			CAData:       readKey("D:/vagrant/.minikube/ca.crt"),
		}
		break
	}

	defer tty.Close()

	tty.Connect()

	c.Status(http.StatusOK)
}

func main() {
	server := gin.Default()
	server.GET("/doterm", doterm)
	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
