package rpc

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	conn     *websocket.Conn
	methods  *Methods
	closer   ConnectionMux
	protocol *Protocol
}

func NewWebsocket(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	return conn, nil
}
func NewConnection(conn *websocket.Conn, methods *Methods, closer ConnectionMux) *Connection {
	c := &Connection{conn: conn, methods: methods, closer: closer}
	c.protocol = NewProtocol(c, methods)
	conn.SetReadLimit(maxMessageSize)
	conn.SetPongHandler(c.HandlePong)
	go c.ReadPump()
	return c
}

func (c *Connection) Close() {
	// call close on hub with connection
	if c.closer != nil {
		c.closer.RemoveConnection(c)
	}
	c.SendClose()
	err := c.conn.Close()
	if err != nil {
		log.Printf("error: %v", err)
	}
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
		c.protocol.HandleMessage(msg)
	}
}

func (c *Connection) SendCall(method string, params []any) (any, error) {
	return c.protocol.Call(method, params)
}

func (c *Connection) SendNotify(method string, params []any) error {
	return c.protocol.Notify(method, params)
}

func (c *Connection) HandlePong(string) error {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

func (c *Connection) SendMessage(msg *RpcMessage) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := c.conn.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		if c.closer != nil {
			// close on hub
			c.closer.RemoveConnection(c)
		}
		c.conn.Close()
		return err
	}
	return nil
}

// SendPing sends a ping message to the server.
func (c *Connection) SendPing() error {
	return c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
}

// SendPong sends a pong message to the server.
func (c *Connection) SendPong() error {
	return c.conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
}

// SendClose sends a close message to the server.
func (c *Connection) SendClose() error {
	return c.conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
}

// func (c *RpcClient) writePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Close()
// 	}()
// 	for {
// 		select {
// 		case msg, ok := <-c.send:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				log.Printf("the hub closed the connection")
// 				c.Close()
// 				return
// 			}

// 			err := c.conn.WriteJSON(msg)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				c.Close()
// 				return
// 			}
// 		case <-ticker.C:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
