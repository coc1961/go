package jms

import (
	"time"
)

/***********************************************************************************************
 ** Server Side
 ***********************************************************************************************/

//Server jms server
type Server struct {
	conn  *Connection
	queue string
}

// Send send message
func (s *Server) Send(msg []byte) error {
	// Envio Msg
	err := s.conn.Send(s.queue, msg)
	if err != nil {
		return err
	}
	return nil
}

//Disconnect Disconnect
func (s *Server) Disconnect() {
	tmp := s.conn
	s.conn = nil
	time.Sleep(time.Millisecond * 100)
	tmp.Disconnect()
}

//NewServer create a jms server object
func NewServer(url, user, password, queue string, listener func(*MessageAck)) (*Server, error) {
	var err error
	// Connecto
	conn, err := Connect(url, user, password)
	if err != nil {
		return nil, err
	}

	// Suscribo Ack
	err = conn.SuscribeAck(queue)
	if err != nil {
		return nil, err
	}

	serv := &Server{conn, queue}

	go func() {

		for serv.conn != nil {
			// Leo Ack
			msg, err := conn.ReadAck()
			if err == nil {
				listener(msg)
			}
		}
	}()

	return serv, nil
}

/***********************************************************************************************
 ** Client Side
 ***********************************************************************************************/

//Client jms client
type Client struct {
	conn *Connection
}

//Disconnect Disconnect
func (s *Client) Disconnect() {
	tmp := s.conn
	s.conn = nil
	tmp.Disconnect()
}

//NewClient create a jms client object
func NewClient(url, user, password, queue string, listener func(msg *Message) []byte) (*Client, error) {
	var err error
	// Connecto
	conn, err := Connect(url, user, password)
	if err != nil {
		return nil, err
	}

	conn.Suscribe(queue)

	client := &Client{conn}

	go func() {
		for client.conn != nil {
			// Leo Ack
			msg, err := conn.Read()
			if err == nil {
				resp := listener(msg)
				if resp != nil {
					msg.SendAck(queue, resp)
				}
			}
		}
	}()

	return client, nil
}
