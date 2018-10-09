package client_server_test_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/coc1961/go/jms"
)

func TestServer(t *testing.T) {

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go server(wg)

	wg.Add(1)
	go client(wg)

	wg.Wait()

}

func server(wg *sync.WaitGroup) {
	server, _ := jms.NewServer("localhost:61613", "admin", "admin", "test", func(msg []byte) {
		fmt.Println("Ack =", string(msg))
	})
	server.Send([]byte("Message 1"))
	time.Sleep(time.Second)
	server.Disconnect()
	wg.Done()
}

func client(wg *sync.WaitGroup) {
	client, _ := jms.NewClient("localhost:61613", "admin", "admin", "test", func(msg *jms.Message) []byte {
		fmt.Println("Msg =", string(msg.Message()))
		return msg.Message()
	})
	time.Sleep(time.Second)
	client.Disconnect()
	wg.Done()
}
