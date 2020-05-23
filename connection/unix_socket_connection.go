package connection

import (
	"bufio"
	"errors"
	"net"
)

type UnixSocketConnection struct {
	socketPath  string
	client      net.Conn
	dataHandler DataHandler
}

func NewUnixSocketConnection(socketPath string) *UnixSocketConnection {
	return &UnixSocketConnection{socketPath: socketPath}
}

func (connection *UnixSocketConnection) SetupConnection() error {
	client, err := net.Dial("unix", connection.socketPath)
	if err != nil {
		return err
	}
	connection.client = client

	go connection.readLoop()

	return nil
}

func (connection *UnixSocketConnection) CloseConnection() error {
	if connection.client == nil {
		return errors.New("client isn't initialized yet. Did you forget to call SetupConnection()")
	}

	err := connection.client.Close()
	if err != nil {
		return err
	}

	return nil
}

func (connection *UnixSocketConnection) Send(data []byte) error {
	if connection.client == nil {
		return errors.New("client isn't initialized yet. Did you forget to call SetupConnection()")
	}

	_, err := connection.client.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (connection *UnixSocketConnection) SetOnDataHandler(dataHandler DataHandler) {
	connection.dataHandler = dataHandler
}

func (connection *UnixSocketConnection) readLoop() {
	reader := bufio.NewReader(connection.client)
	for {
		line, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			return
		}
		go connection.onData(line)
	}
}

func (connection *UnixSocketConnection) onData(data []byte) {
	if connection.dataHandler != nil {
		connection.dataHandler(data)
	}
}
