package connection

type DataHandler func([]byte)

type BaseConnection interface {
	SetupConnection() error
	CloseConnection() error
	Send([]byte) error
	SetOnDataHandler(DataHandler)
}
