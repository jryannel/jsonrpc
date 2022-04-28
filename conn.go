package jsonrpc

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	*Protocol
	conn   *websocket.Conn
	closer ConnectionMux
	send   chan *RpcMessage
}

func NewWebSocket(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	conn.SetReadLimit(maxMessageSize)
	conn.SetPongHandler(func(appData string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	return conn, nil
}
func NewConnection(conn *websocket.Conn, methods *Methods, closer ConnectionMux) *Connection {
	c := &Connection{
		conn:   conn,
		closer: closer,
		send:   make(chan *RpcMessage),
	}
	c.Protocol = NewProtocol(c, methods)
	go c.ReadPump()
	go c.writePump()
	return c
}

func (c *Connection) Close() {
	// call close on hub with connection
	if c.closer != nil {
		c.closer.RemoveConnection(c)
	}
	err := c.conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
	if err != nil {
		log.Printf("error: %v", err)
	}
	err = c.conn.Close()
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func (c *Connection) SendMessage(msg *RpcMessage) error {
	c.send <- msg
	return nil
}

func (c *Connection) ReadPump() {
	defer func() {
		c.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		msg := &RpcMessage{}
		err := c.conn.ReadJSON(msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.handleMessage(msg)
	}
}

func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Printf("the hub closed the connection")
				c.Close()
				return
			}
			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				c.Close()
				return
			}
		case <-ticker.C:
			err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				log.Printf("error: %v", err)
				c.Close()
				return
			}
		}
	}
}
