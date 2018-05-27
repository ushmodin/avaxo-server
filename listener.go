package avaxo

import (
	"fmt"
	"log"
	"net"
)

type Listener struct {
	Opts           ForwardOpts
	clientListener net.Listener
	agentListener  net.Listener
}

func NewListener(opts ForwardOpts) (*Listener, error) {
	cl, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.ClientPort))
	if err != nil {
		return nil, err
	}
	log.Printf("Client listener on port %d started\n", opts.ClientPort)

	al, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.AgentPort))
	if err != nil {
		defer cl.Close()
		return nil, err
	}
	log.Printf("Agent listener on port %d started\n", opts.AgentPort)

	return &Listener{Opts: opts, clientListener: cl, agentListener: al}, nil
}

func (self *Listener) NewConnection() *Connection {
	return &Connection{listener: self}
}

func (self *Listener) Close() {
	self.agentListener.Close()
	self.clientListener.Close()
}
