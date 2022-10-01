package connector

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerTty struct {
	*Websocket
	Host        string
	ContainerID string

	tty *types.HijackedResponse
}

func (tty *DockerTty) Connect() error {

	cli, err := client.NewClientWithOpts(client.WithHost(tty.Host), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	exec, err := cli.ContainerExecCreate(tty.ctx, tty.ContainerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		return err
	}

	attach, err := cli.ContainerExecAttach(tty.ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return err
	}
	tty.tty = &attach

	tty.readerToWebsocket(attach.Conn)
	tty.websocketToWriter(attach.Conn)

	<-tty.ctx.Done()

	return nil
}

func (tty *DockerTty) Close() {

	if tty.tty != nil {
		tty.tty.Close()
	}
}
