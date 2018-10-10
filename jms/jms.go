package jms

import (
	"errors"
	"fmt"

	"github.com/go-stomp/stomp"
	"github.com/go-stomp/stomp/frame"
)

// Connection jms Object
type Connection struct {
	conn    *stomp.Conn
	subs    *stomp.Subscription
	ackSubs *stomp.Subscription
}

// Message jms message Object
type Message struct {
	msg         []byte
	smsg        *stomp.Message
	jms         *Connection
	destination string
	contentType string
}

// Message get message body
func (m *Message) Message() []byte {
	return m.msg
}

// Destination get Destination
func (m *Message) Destination() string {
	return m.destination
}

// ContentType get ContentType
func (m *Message) ContentType() string {
	return m.contentType
}

// MessageAck jms message Object
type MessageAck struct {
	msg         []byte
	destination string
	contentType string
}

// Message get message body
func (m *MessageAck) Message() []byte {
	return m.msg
}

// Destination get Destination
func (m *MessageAck) Destination() string {
	return m.destination
}

// ContentType get ContentType
func (m *MessageAck) ContentType() string {
	return m.contentType
}

var logEnable = false

// Connect Connect to Queue
func Connect(url, user, password string) (*Connection, error) {
	localLog("Connect")
	conn, err := stomp.Dial("tcp", url, stomp.ConnOpt.Login(user, password))
	if err != nil {
		return nil, err
	}
	jms := Connection{conn, nil, nil}
	return &jms, nil
}

// Disconnect Unsuscribe and Disconnect
func (j *Connection) Disconnect() {
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
func (j *Connection) Suscribe(queue string) error {
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
func (j *Connection) Read() (*Message, error) {
	localLog("Read")
	if j.subs == nil {
		return nil, errors.New("Invalid Subscription")
	}
	msg, err := j.subs.Read()
	if err != nil {
		return nil, err
	}
	return &Message{msg.Body, msg, j, msg.Destination, msg.ContentType}, nil
}

// Send Send message
func (j *Connection) Send(queue string, msg []byte) error {
	localLog("Send")
	return j.conn.Send(queue, "", msg, stomp.SendOpt.Receipt)
}

// SuscribeAck Suscribe Queue
func (j *Connection) SuscribeAck(queue string) error {
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
func (j *Connection) ReadAck() (*MessageAck, error) {
	localLog("ReadAck")
	if j.ackSubs == nil {
		return nil, errors.New("Invalid Subscription")
	}
	msg, err := j.ackSubs.Read()
	if err != nil {
		return nil, err
	}
	return &MessageAck{msg.Body, msg.Destination, msg.ContentType}, nil
}

// SendAck Send ack to sender
func (m *Message) SendAck(queue string, msg []byte) error {
	localLog("SendAck")
	err := m.jms.conn.Send(queue+"_ack", "", msg, func(*frame.Frame) error { return nil })
	if err == nil {
		m.jms.conn.Ack(m.smsg)
	}
	return err
}

// SendNack Send ack to sender
func (m *Message) SendNack() error {
	localLog("SendNack")
	err := m.jms.conn.Nack(m.smsg)
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
