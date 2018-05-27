package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ushmodin/avaxo-server"
)

func main() {
	cfg, err := avaxo.NewServerConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	mgr, err := avaxo.NewManager(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err = mgr.Greetings()
			if err != nil {
				log.Printf("Error while listent greetings %s", err)
			}
			<-time.After(30 * time.Second)
		}
	}()

	srv, err := avaxo.NewForwardServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	for {
		err = srv.CreateListeners()
		if err != nil {
			log.Fatal(err)
		}
		<-time.After(60 * time.Second)
	}
	fmt.Println("vim-go")
}
