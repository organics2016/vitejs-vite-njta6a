package connector

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerTty struct {
	Websocket
	Host        string
	ContainerID string

	tty types.HijackedResponse
}

func (dockerTty *DockerTty) Connect() {

	dockerTty.initWebSocket()
	defer dockerTty.Close()

	cli, err := client.NewClientWithOpts(client.WithHost(dockerTty.Host), client.WithAPIVersionNegotiation())
	if err != nil {
		dockerTty.outputError(err)
		return
	}

	exec, err := cli.ContainerExecCreate(dockerTty.ctx, dockerTty.ContainerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		dockerTty.outputError(err)
		return
	}

	tty, err := cli.ContainerExecAttach(dockerTty.ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		dockerTty.outputError(err)
		return
	}
	dockerTty.tty = tty

	dockerTty.readerToWebsocket(tty.Conn)
	dockerTty.websocketToWriter(tty.Conn)

}

func (dockerTty *DockerTty) Close() {
	<-dockerTty.ctx.Done()

	if dockerTty.tty != (types.HijackedResponse{}) {
		dockerTty.tty.Close()
	}
	if dockerTty.websocket != nil {
		dockerTty.websocket.Close()
	}
}
