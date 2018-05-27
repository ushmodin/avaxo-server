package avaxo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomPort() int {
	return int(rnd.Int63n(5000) + 5000)
}

func TestClose(t *testing.T) {
	opts := ForwardOpts{TargetHost: "localhost", TargetPort: randomPort(), ClientPort: randomPort(), AgentPort: randomPort()}
	var agent Agent
	app := NewForwarder(agent)
	listener, err := app.Start(opts)
	if err != nil {
		t.Fatalf("Can't start listener %s", err)
	}
	go app.Loop(listener)
	listener.Close()
	listener, err = app.Start(opts)
	if err != nil {
		t.Fatalf("Can't start listener %s", err)
	}
	go app.Loop(listener)
	listener.Close()
}

func TestForwardEchoServer(t *testing.T) {
	opts := ForwardOpts{TargetHost: "localhost", TargetPort: randomPort(), ClientPort: randomPort(), AgentPort: randomPort()}

	echo, err := NewEchoServer(opts.TargetPort)
	if err != nil {
		t.Fatalf("Can't run echo server %s", err)
	}
	defer echo.Close()
	go echo.Start()

	var agent Agent
	app := NewForwarder(agent)
	listener, err := app.Start(opts)
	if err != nil {
		t.Fatalf("Can't start listener %s", err)
	}
	go app.Loop(listener)

	con, err := net.Dial("tcp", fmt.Sprintf(":%d", opts.TargetPort))
	if err != nil {
		t.Fatalf("Can't start client connection %s", err)
	}
	var returned bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		io.Copy(&returned, con)
		wg.Done()
		log.Println("Input finished")
	}()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var testdata bytes.Buffer

	io.CopyN(&testdata, rnd, rnd.Int63n(10000))
	io.Copy(con, bytes.NewBuffer(testdata.Bytes()))
	<-time.After(1 * time.Second)
	con.Close()

	wg.Wait()
	if bytes.Compare(returned.Bytes(), testdata.Bytes()) != 0 {
		t.Fatalf("Illegal received data %d %d", returned.Len(), testdata.Len())
	}
	listener.Close()
}

func TestConnectionTimeout(t *testing.T) {
	t.Skip()
	opts := ForwardOpts{TargetHost: "localhost", TargetPort: randomPort(), ClientPort: randomPort(), AgentPort: randomPort()}

	var agent EmptyAgent
	app := NewForwarder(agent)
	listener, err := app.Start(opts)
	if err != nil {
		t.Fatalf("Can't start listener %s", err)
	}
	go app.Loop(listener)

	ch := make(chan net.Conn)
	go func() {
		con, err := net.Dial("tcp", ":5000")
		if err != nil {
			t.Fatalf("Can't start client connection %s", err)
		}
		ch <- con
	}()

	select {
	case con := <-ch:
		_, err := con.Read(make([]byte, 1))
		if err != io.EOF {
			t.Fatalf("Connection is not closed %s", err)
		}
	case <-time.After(40 * time.Second):
		t.Fatal("Connectio close timeout")

	}
	listener.Close()
	if len(app.connections) > 0 {
		t.Fatalf("Connections not closed")
	}
}
