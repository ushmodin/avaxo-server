package avaxo

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type Db struct {
	connection string
}

func NewDb(connection string) *Db {
	return &Db{connection: connection}
}

func (self *Db) AllForwardOpts() ([]ForwardOpts, error) {
	db, err := sql.Open("postgres", self.connection)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select id, agent_port, client_port, target_host, target_port, host_id from forward")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var res []ForwardOpts
	for rows.Next() {
		var agentPort, clientPort, targetPort int
		var hostID, id int64
		var targetHost string

		err = rows.Scan(&id, &agentPort, &clientPort, &targetHost, &targetPort, &hostID)
		if err != nil {
			return nil, err
		}
		res = append(res, ForwardOpts{Id: id, AgentPort: agentPort, ClientPort: clientPort, TargetHost: targetHost, TargetPort: targetPort, AgentId: hostID})
	}
	return res, nil
}

func (self *Db) getOrCreateHost(name string) (int64, error) {
	db, err := sql.Open("postgres", self.connection)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	rows, err := db.Query("select id from host where common_name = $1", name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var id int64
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	} else {
		log.Printf("Agent %s not found", name)
		rows, err = db.Query("insert into host (id, title, common_name) values (nextval('hibernate_sequence'), $1, $2) returning id",
			"Unknown host",
			name)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&id)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("Can't fetch sql")
		}
	}
	return id, nil
}
