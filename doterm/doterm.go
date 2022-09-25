package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
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

func ping(c *gin.Context) {

	//前置 : 用户管理，

	// 1. 携带user token 和 目标容器私钥
	// 2. 验证user token
	// 3. 通过目标容器token拿到连接config
	// 4. 通过连接config创建连接，每个websocket对应一个ssh(页面刷新后重新创建ssh)

	//cc := &connector.SSHTty{
	//	Websocket: connector.Websocket{
	//		Request:  c.Request,
	//		Response: c.Writer},
	//
	//	Host:      "127.0.0.1",
	//	Username:  "vagrant",
	//	SecretKey: readKey("D:/vagrant/.vagrant/machines/default/virtualbox/private_key"),
	//	Port:      2222,
	//}

	//cc := &connector.DockerTty{
	//	Websocket: connector.Websocket{
	//		Request:  c.Request,
	//		Response: c.Writer},
	//
	//	Host:        "tcp://127.0.0.1:2375",
	//	ContainerID: "62c41d9cf865b22ba5de8e45462b5744ae34ffd056dbab48542ff1e48c690678",
	//}

	cc := &connector.K8STty{
		Websocket: connector.Websocket{
			Request:  c.Request,
			Response: c.Writer},

		Host:         "https://127.0.0.1:49154",
		PodNamespace: "default",
		PodName:      "shell-demo",
		CertData:     readKey("D:/vagrant/.minikube/profiles/multinode-demo/client.crt"),
		KeyData:      readKey("D:/vagrant/.minikube/profiles/multinode-demo/client.key"),
		CAData:       readKey("D:/vagrant/.minikube/ca.crt"),
	}

	//if err := cc.Connect(); err != nil {
	//	err := err.(stackTracer)
	//	fmt.Printf("%+v\n", err.StackTrace()[0:]) // top two frames
	//	c.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	cc.Connect()
	c.Status(http.StatusOK)

	//st := err.StackTrace()
	//fmt.Printf("%+v", st[0:]) // top two frames

	//if ok {
	//	panic("oops, err does not implement stackTracer")
	//}

}

func main() {
	server := gin.Default()
	server.GET("/ping", ping)
	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
