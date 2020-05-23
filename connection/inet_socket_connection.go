package connection

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type InetSocketConnection struct {
	bindAddr string
	port uint16
	client      net.Conn
	dataHandler DataHandler
}

func NewInetSocketConnection(host string, port uint16) *InetSocketConnection {
	return &InetSocketConnection{bindAddr: host, port: port}
}

func (connection *InetSocketConnection) SetupConnection() error {
	client, err := net.Dial("tcp", fmt.Sprintf("%s:%d", connection.bindAddr, connection.port))
	if err != nil {
		return err
	}
	connection.client = client

	go connection.readLoop()

	return nil
}

func (connection *InetSocketConnection) CloseConnection() error {
	if connection.client == nil {
		return errors.New("client isn't initialized yet. Did you forget to call SetupConnection()")
	}

	err := connection.client.Close()
	if err != nil {
		return err
	}

	return nil
}

func (connection *InetSocketConnection) Send(data []byte) error {
	if connection.client == nil {
		return errors.New("client isn't initialized yet. Did you forget to call SetupConnection()")
	}

	_, err := connection.client.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (connection *InetSocketConnection) SetOnDataHandler(dataHandler DataHandler) {
	connection.dataHandler = dataHandler
}

func (connection *InetSocketConnection) readLoop() {
	reader := bufio.NewReader(connection.client)
	for {
		line, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			return
		}
		go connection.onData(line)
	}
}

func (connection *InetSocketConnection) onData(data []byte) {
	if connection.dataHandler != nil {
		connection.dataHandler(data)
	}
}

