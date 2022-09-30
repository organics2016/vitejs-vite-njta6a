package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"organics.ink/doterm/cloud"
	"organics.ink/doterm/connector"
	"os"
)

func readKey(filePath string) []byte {
	key, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return key
}

type Connection struct {
	Token string `form:"token" binding:"required"`
	Param string `form:"param"`
}

func doterm(c *gin.Context) {

	websocket := connector.InitWebSocket(c.Writer, c.Request, context.Background())

	var connection Connection
	if err := c.ShouldBindQuery(&connection); err != nil {
		websocket.OutputError(err)
		c.Status(http.StatusBadRequest)
		return
	}

	// TODO 获取 连接中心 服务配置
	// 使用token请求连接中心
	// 获取连接信息，开始连接

	URL, err := url.Parse(BootOptions.Server)
	if err != nil {
		panic(err)
	}
	params := url.Values{}
	params.Set("token", connection.Token)
	params.Set("param", connection.Param)
	URL.RawQuery = params.Encode()

	request, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		return
	}
	request.Header.Add(restful.HEADER_ContentType, "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	connData := &cloud.ConnData{}
	json.Unmarshal(body, connData)

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

type bootOptions struct {
	Server string
}

var BootOptions = &bootOptions{}

func parseFlags() {
	flag.StringVar(&BootOptions.Server, "s", "http://127.0.0.1/", "aaaaaa")
	flag.Parse()
}

func main() {
	parseFlags()
	fmt.Printf("%v+", BootOptions)

	server := gin.Default()
	server.GET("/doterm", doterm)
	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
