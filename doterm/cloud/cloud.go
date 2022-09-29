package cloud

type Authorization struct {
	Api string
}

type ConnData struct {
	Host string
	Port int
	Type string

	// host
	Username string
	Password string
	PubKey   string
	PriKey   string

	// docker
	ContainerID string

	// k8s
	PodNamespace string
	PodName      string

	// docker&k8s
	CertData string
	KeyData  string
	CAData   string
}

type Cloud interface {
	Authorize() ConnData
}
