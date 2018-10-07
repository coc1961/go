package jms_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/coc1961/go/jms"
)

//QUEUENAME Queue Name
const QUEUENAME = "testQueue"

//CONT Number of Message to process
var CONT = 100

var sent = 0
var recv = 0

func TestConnect(t *testing.T) {

	jms.SetLogEnable(false)

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
	var err error

	// Connecto
	conn, err := jms.Connect("localhost:61613", "admin", "admin")
	printError("server", err)

	// Suscribo Ack
	err = conn.SuscribeAck(QUEUENAME)
	printError("server", err)

	for i := 0; i < CONT; i++ {

		// Envio Msg
		err = conn.Send(QUEUENAME, []byte(fmt.Sprintf("Message %d", i+1)))
		printError("server", err)

		// Leo Ack
		msg, err := conn.ReadAck()
		printError("server", err)
		fmt.Println("Ack =", string(msg))

		sent++

	}

	// Desconecto
	conn.Disconnect()
	wg.Done()
}

func client(wg *sync.WaitGroup) {

	// Connecto
	conn, err := jms.Connect("localhost:61613", "admin", "admin")
	printError("client", err)

	var cont = 0
	// Suscribo a la cola con un listener
	go conn.SuscribeListener(QUEUENAME, func(msg *jms.Message) []byte {
		fmt.Println("Msg =", string(msg.Msg))
		cont++
		recv++
		return []byte(string(msg.Msg) + ".Ok")
	})

	// Espero que se procesen los mensajes
	for cont < CONT {
		time.Sleep(time.Millisecond * 100)
	}

	// Desconecto
	conn.Disconnect()
	wg.Done()
}

func printError(quien string, err error) {
	if err != nil {
		fmt.Println(quien, err)
	}
}
