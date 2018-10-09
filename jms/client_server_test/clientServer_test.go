package client_server_test_test

import (
	"fmt"
	"sync"
	"testing"

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
	//wg.Done()
}

func client(wg *sync.WaitGroup) {
	jms.NewClient("localhost:61613", "admin", "admin", "test", func(msg *jms.Message) []byte {
		fmt.Println("Msg =", string(msg.Message()))
		return msg.Message()
	})

	//wg.Done()
}
