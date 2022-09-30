package main

import (
	"encoding/json"
	"github.com/emicklei/go-restful/v3"
	"io"
	"net/http"
	"net/url"
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

type ConnData struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
	Type string `json:"type,omitempty"`

	// host
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	PubKey   string `json:"pubKey,omitempty"`
	PriKey   string `json:"priKey,omitempty"`

	// docker
	ContainerID string `json:"containerID,omitempty"`

	// k8s
	PodNamespace string `json:"podNamespace,omitempty"`
	PodName      string `json:"podName,omitempty"`

	// docker&k8s
	CertData string `json:"certData,omitempty"`
	KeyData  string `json:"keyData,omitempty"`
	CAData   string `json:"CAData,omitempty"`
}

func getConnData(cp *ConnParam) (*ConnData, error) {
	URL, err := url.Parse(bootOptions.Server)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("token", cp.Token)
	params.Set("param", cp.Param)
	URL.RawQuery = params.Encode()

	request, err := http.NewRequest(http.MethodGet, URL.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add(restful.HEADER_ContentType, "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	connData := &ConnData{}
	if err := json.Unmarshal(body, connData); err != nil {
		return nil, err
	}
	return connData, nil
}

func connTTY(connData *ConnData, websocket *connector.Websocket) {
	var tty connector.TTY

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
}

func conn(cp *ConnParam, websocket *connector.Websocket) error {

	connData, err := getConnData(cp)
	if err != nil {
		return err
	}

	connTTY(connData, websocket)

	return nil
}
