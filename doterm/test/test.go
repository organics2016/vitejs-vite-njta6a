package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"os"
)

func test() {
	bs_UTF16LE, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(), []byte("1"))
	bs_UTF16BE, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder(), []byte("1"))

	bs_UTF16LEN, _ := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder().Bytes([]byte("1"))
	bs_UTF16BEN, _ := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder().Bytes([]byte("1"))

	bs_UTF8LE, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), bs_UTF16LE)
	bs_UTF8BE, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), bs_UTF16BE)

	fmt.Printf("%v\n%v\n%v\n%v\n%v\n%v\n", bs_UTF16LE, bs_UTF16BE, bs_UTF16LEN, bs_UTF16BEN, bs_UTF8LE, bs_UTF8BE)
}

func test2() error {

	return errors.New("dddddd")
}

func test3() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("http://127.0.0.1:2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}

func test4() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("http://127.0.0.1:2376"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}

}

func test5() string {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("http://127.0.0.1:2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	response, err := cli.ContainerExecCreate(ctx,
		"62c41d9cf865b22ba5de8e45462b5744ae34ffd056dbab48542ff1e48c690678",
		types.ExecConfig{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			Cmd:          []string{"/bin/bash"},
		})
	if err != nil {
		panic(err)
	}

	return response.ID
}

func test7() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://127.0.0.1:2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	id := test5()
	println(id)
	response, err := cli.ContainerExecAttach(ctx, id, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		panic(err)
	}

	// 关闭I/O
	//defer response.Close()
	// 输入
	response.Conn.Write([]byte("ls\r"))
	// 输出
	scanner := bufio.NewScanner(response.Conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}

func test8() {

	if true {
		panic("bbb")
	}

	defer print("aaaaa")

}

func main() {
	test8()
	println("dddddd")
}
