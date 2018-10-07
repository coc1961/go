package jms

import (
	"errors"
	"fmt"

	"github.com/go-stomp/stomp"
)

// JMS jms Object
type JMS struct {
	conn    *stomp.Conn
	subs    *stomp.Subscription
	ackSubs *stomp.Subscription
}

// Message jms message Object
type Message struct {
	Msg  []byte
	smsg *stomp.Message
	jms  *JMS
}

var logEnable = false

// Connect Connect to Queue
func Connect(url, user, password string) (*JMS, error) {
	localLog("Connect")
	conn, err := stomp.Dial("tcp", url, stomp.ConnOpt.Login(user, password))
	if err != nil {
		return nil, err
	}
	jms := JMS{conn, nil, nil}
	return &jms, nil
}

// Disconnect Unsuscribe and Disconnect
func (j *JMS) Disconnect() {
	localLog("Disconnect")
	if j.subs != nil {
		j.subs.Unsubscribe()
		j.subs = nil
	}
	if j.ackSubs != nil {
		j.ackSubs.Unsubscribe()
		j.ackSubs = nil
	}
	if j.conn != nil {
		j.conn.Disconnect()
		j.conn = nil
	}
}

// Suscribe Suscribe Queue
func (j *JMS) Suscribe(queue string) error {
	localLog("Suscribe")
	if queue == "" {
		return errors.New("Invalid Queue")
	}

	sub, err := j.conn.Subscribe(queue, stomp.AckClientIndividual)
	if err != nil {
		return err
	}

	j.subs = sub

	return nil
}

// Read Read Message
func (j *JMS) Read() (*Message, error) {
	localLog("Read")
	if j.subs == nil {
		return nil, errors.New("Invalid Subscription")
	}
	msg, err := j.subs.Read()
	if err != nil {
		return nil, err
	}
	return &Message{msg.Body, msg, j}, nil
}

// Send Send message
func (j *JMS) Send(queue string, msg []byte) error {
	localLog("Send")
	return j.conn.Send(queue, "", msg, stomp.SendOpt.Receipt)
}

// SuscribeAck Suscribe Queue
func (j *JMS) SuscribeAck(queue string) error {
	localLog("SuscribeAck")
	if queue == "" {
		return errors.New("Invalid Queue")
	}
	sub, err := j.conn.Subscribe(queue+"_ack", stomp.AckAuto)
	if err != nil {
		return err
	}
	j.ackSubs = sub
	return nil
}

// ReadAck Read Ack
func (j *JMS) ReadAck() ([]byte, error) {
	localLog("ReadAck")
	if j.ackSubs == nil {
		return nil, errors.New("Invalid Subscription")
	}
	msg, err := j.ackSubs.Read()
	if err != nil {
		return nil, err
	}
	return msg.Body, nil
}

// SendAck Send ack to sender
func (j *Message) SendAck(queue string, msg []byte) error {
	localLog("SendAck")
	err := j.jms.conn.Send(queue+"_ack", "", msg, stomp.SendOpt.Receipt)
	if err == nil {
		j.jms.conn.Ack(j.smsg)
	}
	return err
}

// SetLogEnable enable internal log
func SetLogEnable(b bool) {
	logEnable = b
}

func localLog(msg string) {
	if logEnable {
		fmt.Println(msg)
	}

}
