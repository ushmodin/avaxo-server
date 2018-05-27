package avaxo

import (
	"fmt"
	"io"
	"log"
	"net"
)

type EchoServer struct {
	listener net.Listener
}

type Client struct {
	conn net.Conn
}

func NewEchoServer(port int) (*EchoServer, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return &EchoServer{listener: l}, nil
}

func (self *EchoServer) Start() {
	defer self.Close()
	for {
		client, err := self.accept()
		if err != nil {
			log.Printf("Can't accept %s", err)
			continue
		}
		go func() {
			defer client.Close()
			client.Forward()
		}()
	}

}

func (self *EchoServer) accept() (*Client, error) {
	client, err := self.listener.Accept()
	if err != nil {
		return nil, err
	}
	return &Client{conn: client}, nil
}

func (self *EchoServer) Close() {
	self.listener.Close()
}

func (self *Client) Forward() {
	io.Copy(self.conn, self.conn)
}

func (self *Client) Close() {
	self.conn.Close()
}
