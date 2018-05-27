package avaxo

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ServerConfig struct {
	RabbitConnectionString string `json:"rabbitConnectionString"`
	DbConnectionString     string `json:"dbConnectionString"`
	ManagerHost            string
}

func NewServerConfig(file string) (*ServerConfig, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg ServerConfig
	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type ForwardServer struct {
	config        *ServerConfig
	db            *Db
	forwarder     *Forwarder
	agentNotifier *AgentNotifier
	listeners     map[int64]*Listener
}

func NewForwardServer(cfg *ServerConfig) (*ForwardServer, error) {

	agentNotifier := AgentNotifier(NewAmqpAgent(cfg.RabbitConnectionString, cfg.ManagerHost))

	return &ForwardServer{
		config:        cfg,
		agentNotifier: &agentNotifier,
		db:            NewDb(cfg.DbConnectionString),
		forwarder:     NewForwarder(agentNotifier),
		listeners:     make(map[int64]*Listener),
	}, nil
}

func (self *ForwardServer) CreateListeners() error {
	log.Println("Refresh listeners")
	opts, err := self.db.AllForwardOpts()
	if err != nil {
		return err
	}
	for _, item := range opts {
		_, ok := self.listeners[item.Id]
		if ok {
			continue
		}
		lstr, err := self.forwarder.Start(item)
		if err == nil {
			self.listeners[item.Id] = lstr
			go self.forwarder.Loop(lstr)
		} else {
			log.Print(err)
		}
	}
	return nil
}
