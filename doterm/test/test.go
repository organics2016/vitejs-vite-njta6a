package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"time"
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

func test9() {
	//var kubeconfig *string
	//if home := homedir.HomeDir(); home != "" {
	//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()

	kubeconfig := "D:/vagrant/config"

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	config = &rest.Config{
		Host: "https://127.0.0.1:49154",
		TLSClientConfig: rest.TLSClientConfig{
			CertFile: "D:/vagrant/.minikube/profiles/multinode-demo/client.crt",
			KeyFile:  "D:/vagrant/.minikube/profiles/multinode-demo/client.key",
			CAFile:   "D:/vagrant/.minikube/ca.crt",
		},
	}

	fmt.Printf("%+v\n", config)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	req := clientset.CoreV1().
		RESTClient().
		Get().
		Namespace("default").
		Resource("pods").
		Name("shell-demo").
		SubResource("exec").
		VersionedParams(
			&corev1.PodExecOptions{
				Command: []string{"/bin/sh"},
				Stdin:   true,
				Stdout:  true,
				Stderr:  true,
				TTY:     true,
			}, scheme.ParameterCodec)

	fmt.Println(req)

	exec, err := remotecommand.NewSPDYExecutor(config, "GET", req.URL())
	if err != nil {
		fmt.Println(err)
		return
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		fmt.Println("error streaming connection", err)
		return
	}

	time.Sleep(10 * time.Second)
}

func test10() {

	p := make([]byte, 10, 20)
	p = nil
	fmt.Println(len(p))

}

func test11() {
	file1, err := os.ReadFile("D:/vagrant/authorized_keys")
	if err != nil {
		panic(err)
	}
	file2, err := os.ReadFile("D:/vagrant/.vagrant/machines/default/virtualbox/private_key")
	if err != nil {
		panic(err)
	}

	str1 := base64.StdEncoding.EncodeToString(file1)
	println(str1)
	str2 := base64.StdEncoding.EncodeToString(file2)
	println(str2)

}

func main() {
	test11()
	println("dddddd")
}
