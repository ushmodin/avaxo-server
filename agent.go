package avaxo

import (
	"fmt"
	"io"
	"net"
)

type Agent struct {
}

func (self Agent) Notify(opts ForwardOpts) error {
	tgt, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opts.TargetHost, opts.TargetPort))
	if err != nil {
		return err
	}
	agt, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", opts.AgentPort))
	if err != nil {
		tgt.Close()
		return err
	}
	go func() {
		io.Copy(tgt, agt)
		tgt.Close()
	}()
	go func() {
		io.Copy(agt, tgt)
		agt.Close()
	}()

	return nil
}
