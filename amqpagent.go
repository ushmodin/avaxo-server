package avaxo

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type AmqpAgent struct {
	connectionString string
	managerHost      string
}

func NewAmqpAgent(connectionString, managerHost string) *AmqpAgent {
	return &AmqpAgent{connectionString: connectionString, managerHost: managerHost}
}

func (self *AmqpAgent) Notify(opts ForwardOpts) error {
	var cmd struct {
		Forward struct {
			Target struct {
				Port int    `json:"port"`
				Host string `json:"host"`
			} `json:"target"`
			Manager struct {
				Port int    `json:"port"`
				Host string `json:"host"`
			} `json:"manager"`
		} `json:"forward"`
	}
	cmd.Forward.Target.Host = opts.TargetHost
	cmd.Forward.Target.Port = opts.TargetPort
	cmd.Forward.Manager.Port = opts.AgentPort
	cmd.Forward.Manager.Host = self.managerHost
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return self.sendMessage(fmt.Sprintf("avaxo.agent_%d", opts.AgentId), msg)

}

func (self *AmqpAgent) sendMessage(queue string, msg []byte) error {
	conn, err := amqp.Dial(self.connectionString)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		"",
		queue,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})

	return err

}
