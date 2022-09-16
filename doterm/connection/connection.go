package connection

type ConnectorManager interface {
	Connect() error
}
