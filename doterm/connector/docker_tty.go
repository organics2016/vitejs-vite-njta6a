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

func (tty *DockerTty) Connect() {

	cli, err := client.NewClientWithOpts(client.WithHost(tty.Host), client.WithAPIVersionNegotiation())
	if err != nil {
		tty.OutputError(err)
		return
	}

	exec, err := cli.ContainerExecCreate(tty.ctx, tty.ContainerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		tty.OutputError(err)
		return
	}

	attach, err := cli.ContainerExecAttach(tty.ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		tty.OutputError(err)
		return
	}
	tty.tty = attach

	tty.readerToWebsocket(attach.Conn)
	tty.websocketToWriter(attach.Conn)

}

func (tty *DockerTty) Close() {
	<-tty.ctx.Done()

	if tty.tty != (types.HijackedResponse{}) {
		tty.tty.Close()
	}
	if tty.Websocket.wsConn != nil {
		tty.Websocket.wsConn.Close()
	}
}
