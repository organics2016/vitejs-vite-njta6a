package connector

import (
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

type K8STty struct {
	*Websocket
	Host         string
	PodNamespace string
	PodName      string
	CertData     []byte
	KeyData      []byte
	CAData       []byte
}

func (tty *K8STty) Connect() error {

	config := &rest.Config{
		Host: tty.Host,
		TLSClientConfig: rest.TLSClientConfig{
			CertData: tty.CertData,
			KeyData:  tty.KeyData,
			CAData:   tty.CAData,
		},
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	req := clientSet.CoreV1().
		RESTClient().
		Post().
		Namespace(tty.PodNamespace).
		Resource("pods").
		Name(tty.PodName).
		SubResource("exec").
		VersionedParams(
			&corev1.PodExecOptions{
				Command: []string{"/bin/sh"},
				Stdin:   true,
				Stdout:  true,
				Stderr:  true,
				TTY:     true,
			}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  tty,
		Stdout: tty,
		Stderr: tty,
		Tty:    true,
	})

	tty.cancel()
	<-tty.ctx.Done()

	return nil
}

func (tty *K8STty) Close() {

}

func (tty *K8STty) Read(p []byte) (n int, err error) {
	_, message, err := tty.Websocket.wsConn.ReadMessage()
	if err != nil {
		tty.cancel()
		return 0, err
	}
	//fmt.Printf("%v : %t\n", string(p), utf8.Valid(p))
	//message, _ = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder().Bytes(message)
	//fmt.Printf("dist: [%+v] size : %d --- src: [%+v] size : %d \n", string(p), len(p), string(message), len(message))
	c := copy(p, message)
	//fmt.Printf("dist: [%+v] size : %d --- src: [%+v] size : %d \n", string(p), len(p), string(message), len(message))
	return c, nil
}

func (tty *K8STty) Write(p []byte) (n int, err error) {
	//fmt.Printf("Write: [%+v] size : %d utf8: %t\n", string(p), len(p), utf8.Valid(p))
	//e, _ := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder().Bytes(p)
	//fmt.Printf("EWrite: [%+v] size : %d utf8: %t\n", string(e), len(e), utf8.Valid(e))
	if err := tty.Websocket.wsConn.WriteMessage(websocket.BinaryMessage, p); err != nil {
		tty.cancel()
		return 0, err
	}
	return len(p), nil
}
