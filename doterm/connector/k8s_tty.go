package connector

type K8STty struct {
	Websocket
	Host      string
	Username  string
	Password  string
	SecretKey []byte
	Port      int
}

func (tty *K8STty) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (tty *K8STty) Write(p []byte) (n int, err error) {
	return 0, nil
}
