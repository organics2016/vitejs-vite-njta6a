package connector

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"net/http"
)

type DockerTTY struct {
	Request     *http.Request
	Writer      http.ResponseWriter
	Host        string
	ContainerID string

	websocket *websocket.Conn
}

func (d *DockerTTY) Connect() error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("http://127.0.0.1:2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	exec, err := cli.ContainerExecCreate(ctx, d.ContainerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		panic(err)
	}

	tty, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		panic(err)
	}

	tty.Conn.Write([]byte("ls\r"))
	//scanner := bufio.NewScanner(response.Conn)
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//}
	return nil

}
