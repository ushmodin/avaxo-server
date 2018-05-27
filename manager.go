package avaxo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	DEFAULT_TTL = int32(120000)
)

type Manager struct {
	config *ServerConfig
	db     *Db
}

func NewManager(cfg *ServerConfig) (*Manager, error) {
	return &Manager{config: cfg, db: NewDb(cfg.DbConnectionString)}, nil
}

func (self *Manager) Greetings() error {
	conn, err := amqp.Dial(self.config.RabbitConnectionString)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	greetingsQueue := "avaxo.greetings"
	err = self.createGreetingQueue(conn, greetingsQueue)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(greetingsQueue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		self.greeting(conn, msg.Body)
	}
	return nil
}

func (mgr *Manager) greeting(conn *amqp.Connection, body []byte) {
	var greeting struct {
		From          string `json:"from"`
		CallbackQueue string `json:"callbackQueue"`
	}
	err := json.Unmarshal(body, &greeting)
	if err != nil {
		log.Printf("Can't unmarshal greeting message %s", err)
	}
	id, err := mgr.db.getOrCreateHost(greeting.From)
	if err != nil {
		log.Printf("Can't get agent id %s", err)
	}

	commandQueue := fmt.Sprintf("avaxo.agent_%d", id)

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Can't create chanel %s", err)
	}
	defer ch.Close()

	err = mgr.createAgentQueue(conn, commandQueue)
	if err != nil {
		log.Printf("Can't declare queue %s", err)
	}

	err = mgr.sendGreetingResponse(ch, greeting.CallbackQueue, commandQueue, id)
}

func (mgr *Manager) createAgentQueue(conn *amqp.Connection, name string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	params := make(amqp.Table)
	params["x-message-ttl"] = DEFAULT_TTL
	_, err = ch.QueueDeclarePassive(name, true, false, false, false, params)
	if err == nil {
		return nil
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(name, true, false, false, false, params)
	return err
}

func (mgr *Manager) createGreetingQueue(conn *amqp.Connection, name string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	params := make(amqp.Table)
	params["x-message-ttl"] = DEFAULT_TTL

	_, err = ch.QueueDeclarePassive(name, true, false, false, false, params)

	if err != nil {
		ch, err := conn.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		_, err = ch.QueueDeclare(name, true, false, false, false, params)
		if err != nil {
			return err
		}
		log.Printf("Queue %s created", name)
	}
	return nil
}

func (*Manager) sendGreetingResponse(ch *amqp.Channel, queue, commandQueue string, id int64) error {
	var response struct {
		ID             int64  `json:"id"`
		CommandQueueu  string `json:"commandQueue"`
		HeartbeatQueue string `json:"heartbeatQueue"`
	}
	response.ID = id
	response.CommandQueueu = commandQueue
	response.HeartbeatQueue = "avaxo.heartbeat"

	responseBody, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return ch.Publish(
		"",
		queue,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        responseBody,
		})
}
