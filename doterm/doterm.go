package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"io"
	"net/http"
	"net/url"
	"organics.ink/doterm/connector"
)

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

func connTTY(connData *ConnData, websocket *connector.Websocket) error {
	var tty connector.TTY

	switch connData.Type {
	case "host":

		pubKey, err := base64.StdEncoding.DecodeString(connData.PubKey)
		if err != nil {
			return err
		}
		priKey, err := base64.StdEncoding.DecodeString(connData.PriKey)
		if err != nil {
			return err
		}

		tty = &connector.SSHTty{
			Websocket: websocket,

			Host:     connData.Host,
			Port:     connData.Port,
			Username: connData.Username,
			Password: connData.Password,
			PubKey:   pubKey,
			PriKey:   priKey,
		}
		break
	case "docker":

		addr := fmt.Sprintf("tcp://%s:%d", connData.Host, connData.Port)
		tty = &connector.DockerTty{
			Websocket: websocket,

			Host:        addr,                 //"tcp://127.0.0.1:2375"
			ContainerID: connData.ContainerID, // "62c41d9cf865b22ba5de8e45462b5744ae34ffd056dbab48542ff1e48c690678"
		}
		break
	case "kubernetes":
		certData, err := base64.StdEncoding.DecodeString(connData.CertData)
		if err != nil {
			return err
		}
		keyData, err := base64.StdEncoding.DecodeString(connData.KeyData)
		if err != nil {
			return err
		}
		caData, err := base64.StdEncoding.DecodeString(connData.CAData)
		if err != nil {
			return err
		}

		addr := fmt.Sprintf("https://%s:%d", connData.Host, connData.Port)
		tty = &connector.K8STty{
			Websocket: websocket,

			Host:         addr,                  //"https://127.0.0.1:49154",
			PodNamespace: connData.PodNamespace, //"default",
			PodName:      connData.PodName,      //"shell-demo",
			CertData:     certData,              //readKey("D:/vagrant/.minikube/profiles/multinode-demo/client.crt"),
			KeyData:      keyData,               //readKey("D:/vagrant/.minikube/profiles/multinode-demo/client.key"),
			CAData:       caData,                //readKey("D:/vagrant/.minikube/ca.crt"),
		}
		break
	}

	defer tty.Close()

	if err := tty.Connect(); err != nil {
		return err
	}

	return nil
}

func conn(cp *ConnParam, websocket *connector.Websocket) error {

	connData, err := getConnData(cp)
	if err != nil {
		return err
	}

	if err := connTTY(connData, websocket); err != nil {
		return err
	}

	return nil
}
