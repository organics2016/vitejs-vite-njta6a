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

	dockerTty.initWebSocket(dockerTty)
	// 如果websocket先断开连接，这里会重复执行一次，当容器先断开连接时或发生意外，在这里释放资源
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

	<-dockerTty.ctx.Done()
}

func (dockerTty *DockerTty) Close() {
	dockerTty.cancel()

	dockerTty.tty.Conn.Write([]byte("exit\r"))
	dockerTty.tty.Close()
	dockerTty.tty.CloseWrite()
}
