package avaxo

import (
	"io"
	"log"
	"net"
)

type Connection struct {
	listener   *Listener
	agentConn  net.Conn
	clientConn net.Conn
}

func (self *Connection) AcceptClient() error {
	con, err := self.listener.clientListener.Accept()
	if err != nil {
		return err
	}
	log.Printf("Client connected on port %d\n", self.listener.Opts.ClientPort)
	self.clientConn = con
	return nil
}

func (self *Connection) AcceptAgent() error {
	con, err := self.listener.agentListener.Accept()
	if err != nil {
		return err
	}
	log.Printf("Agent connected on port %d\n", self.listener.Opts.AgentPort)
	self.agentConn = con
	return nil

}

func (self *Connection) ForwardC2A() {
	log.Printf("Copy traffic from client port %d to agent port %d\n", self.listener.Opts.ClientPort, self.listener.Opts.AgentPort)
	io.Copy(self.agentConn, self.clientConn)
	log.Printf("End Copy traffic from client port %d to agent port %d\n", self.listener.Opts.ClientPort, self.listener.Opts.AgentPort)
	self.agentConn.Close()
}

func (self *Connection) ForwardA2C() {
	log.Printf("Copy traffic from agent port %d to client port %d\n", self.listener.Opts.AgentPort, self.listener.Opts.ClientPort)
	io.Copy(self.clientConn, self.agentConn)
	log.Printf("End copy traffic from agent port %d to client port %d\n", self.listener.Opts.AgentPort, self.listener.Opts.ClientPort)
	self.clientConn.Close()
}

func (self *Connection) Close() {
	if self.agentConn != nil {
		defer self.agentConn.Close()
	}
	defer self.clientConn.Close()
}
