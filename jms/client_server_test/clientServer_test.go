package client_server_test_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/coc1961/go/jms"
)

const CONT = 10000

var sent = 0
var recv = 0

func TestServer(t *testing.T) {

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go server(wg)

	wg.Add(1)
	go client(wg)

	wg.Wait()

	if sent != recv || sent != CONT {
		t.Fatal("Send and Recv Differ")
	}
}

func server(wg *sync.WaitGroup) {
	server, _ := jms.NewServer("localhost:61613", "admin", "admin", "test", func(msg []byte) {
		fmt.Println("Ack =", string(msg))
		sent++
	})
	for i := 0; i < CONT; i++ {
		msg := fmt.Sprintf("Message %d", i)
		server.Send([]byte(msg))
	}
	for sent < CONT {
		time.Sleep(time.Microsecond)
	}
	time.Sleep(time.Second)
	server.Disconnect()
	wg.Done()
}

func client(wg *sync.WaitGroup) {
	client, _ := jms.NewClient("localhost:61613", "admin", "admin", "test", func(msg *jms.Message) []byte {
		fmt.Println("Msg =", string(msg.Message()))
		recv++
		return msg.Message()
	})
	for recv < CONT {
		time.Sleep(time.Microsecond)
	}
	time.Sleep(time.Second)
	client.Disconnect()
	wg.Done()
}
