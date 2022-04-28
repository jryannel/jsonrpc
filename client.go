package jsonrpc

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Max message size allowed from peer.
	maxMessageSize = 512
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type RpcClient struct {
	Conn *Connection
	// Methods *Methods
	Methods
}

func NewRpcClient(conn *websocket.Conn) *RpcClient {
	c := &RpcClient{
		Conn: nil,
		// Methods: NewRegistry(),
		Methods: Methods{
			methods: make(map[string]MethodHandle),
		},
	}
	c.Conn = NewConnection(conn, &c.Methods, nil)
	// go c.readPump()
	// go c.writePump()
	return c
}

// Close closes the connection
func (c *RpcClient) Close() {
	c.Conn.conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeWait))
	c.Conn.Close()
}

// Call sends a request to the server and waits for a response.
func (c *RpcClient) Call(method string, params []any) (any, error) {
	return c.Conn.SendCall(method, params)
}

// Notify sends a notification to the server.
func (c *RpcClient) Notify(method string, params []any) error {
	return c.Conn.SendMessage(MakeNotify(method, params))
}

// Error sends an error message to the server.
func (c *RpcClient) Error(code ErrorCode, message string, data interface{}) error {
	return c.Conn.SendMessage(MakeError(code, message, data))
}

// Send sends a message to the server.
func (c *RpcClient) SendMessage(msg *RpcMessage) error {
	return c.Conn.SendMessage(msg)
}
