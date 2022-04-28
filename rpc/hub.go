package rpc

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ConnectionMux is an interface for a connection multiplexer.
type ConnectionMux interface {
	RemoveConnection(c *Connection)
	BroadcastMessage(msg *RpcMessage)
}

// Hub is the central hub for all connections and method registry.
type Hub struct {
	upgrader websocket.Upgrader
	Connections
	Methods
}

func NewHub() *Hub {
	return &Hub{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Connections: Connections{
			connections: make(map[*Connection]bool),
		},
		Methods: Methods{
			methods: make(map[string]MethodHandle),
		},
	}
}

func (h *Hub) HandleRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := NewConnection(conn, &h.Methods, h)
	h.AddConnection(c)
}

func (h *Hub) Notify(method string, params []any) {
	h.BroadcastMessage(MakeNotify(method, params))
}
