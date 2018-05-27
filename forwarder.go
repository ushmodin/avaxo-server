package avaxo

import (
	"log"
	"sync"
	"time"
)

type ForwardOpts struct {
	Id         int64
	AgentId    int64
	AgentPort  int
	ClientPort int
	TargetHost string
	TargetPort int
}

type Forwarder struct {
	agentNotifier AgentNotifier
	connections   map[*Connection]sync.WaitGroup
}

type AgentNotifier interface {
	Notify(opts ForwardOpts) error
}

func NewForwarder(notifier AgentNotifier) *Forwarder {
	return &Forwarder{agentNotifier: notifier, connections: make(map[*Connection]sync.WaitGroup)}
}

func (self *Forwarder) Start(params ForwardOpts) (*Listener, error) {
	l, err := NewListener(params)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (self *Forwarder) Loop(cl *Listener) {
	for {
		conn := cl.NewConnection()
		err := conn.AcceptClient()
		if err != nil {
			log.Printf("Error while accept client connection %s\n", err)
			return
		}
		go func() {
			ch := make(chan error)
			go func() {
				err := conn.AcceptAgent()
				ch <- err
			}()
			select {
			case err := <-ch:
				if err != nil {
					log.Printf("Error while accept agent connection %s\n", err)
					defer conn.Close()
					return
				}
			case <-time.After(30 * time.Second):
				log.Println("Agent connection timeout")
				defer conn.Close()
				return
			}

			var wg sync.WaitGroup
			self.connections[conn] = wg

			wg.Add(2)
			go func() {
				defer wg.Done()
				conn.ForwardA2C()
			}()
			go func() {
				defer wg.Done()
				conn.ForwardC2A()
			}()

			wg.Wait()
			delete(self.connections, conn)
		}()
		self.agentNotifier.Notify(cl.Opts)
		log.Printf("Send signal to agent")
	}

}
